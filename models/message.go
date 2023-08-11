package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"local_imessage/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 调度
func init() {
	fmt.Println("init...")
	go udpSendProc()
	go udpRecvProc()
}

type Message struct {
	gorm.Model
	UserId     string // 发送者
	TargetId   string // 接收者
	Type       int    // 发送类型: 群聊,私聊,广播等
	Media      int    // 文字,图片,音频
	Content    string // 消息内容
	CreateTime uint64 // 创建时间
	ReadTime   uint64 // 读取时间
	Pic        string // 图片
	Url        string // 链接
	Desc       string // 描述
	Amount     int    // 其它统计
}
type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
	GroupSets     set.Interface   //好友 / 群
}

func (table *Message) TableName() string {
	return "message"
}

// 映射关系
var clientMap = make(map[string]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {

	// 从前端获取的请求参数都是字符串类型的

	// 1.获取参数 以及 检验 token 以及其它合法性
	// token := query.Get("token") //暂时不校验
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")

	isVALID := true // 待完成
	// 升级链接
	conn, err := (&websocket.Upgrader{
		// token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isVALID
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2.初始化 node
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,                       // 这是一个升级后的websocket
		DataQueue:     make(chan []byte, 50),      // 有可能有多个人给一个人发送消息,管道容量暂设定为 50
		GroupSets:     set.New(set.ThreadSafe),    // 线程安全群集合,对其进行写入时,应该是线程安全的
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登录时间
	}
	// 3.获取关系
	// 4.userId 与 node绑定,并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5.多协程发送
	go sendProc(node)
	// 6.多协程完成接收
	go recvProc(node)
	// 仅供测试
	// sendMsg(userId, []byte("欢迎来到聊天室233!"))

	// 将在线用户信息加入到缓存
	SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)

}

// 定向写入User_node_conn
func sendProc(node *Node) {
	fmt.Println("sendProc...")
	// 一直循环等待处理 Node_用户 所发的消息
	for {
		select {
		// 这是一个死循环中的管道,是阻塞式读写的,它可以源源不断地将用户的数据取出来
		case data := <-node.DataQueue:
			// 将用户的消息装进 data
			// 将用户的 data 写进 User_node 中的 Conn
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

// 定向读出User_node_conn
func recvProc(node *Node) {
	fmt.Println("recvProc...")
	for {
		// 将用户的数据源源不断地用 Conn 中读出来
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		// 将消息反序列化到 msg
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println(err)
		}
		// 如果是心跳包,设置心跳开始时间
		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			disPatch(data) // 消息分发
			broadMsg(data) // 将消息广播

			// fmt.Println("[ws] recvProc <<<<< ", string(data))
		}
	}
}

// 将信息全部存进 udpSendChan
var udpSendChan = make(chan []byte, 1024*32)

func broadMsg(data []byte) {
	udpSendChan <- data

}

// 完成数据发送协程,发送到 UDPconn 中
func udpSendProc() {
	fmt.Println("udpSendProc...")
	UDPconn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(10, 30, 0, 159),
		Port: 3000,
	})
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {

		}
	}(UDPconn)
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpSendChan:
			// 非定向发送数据到 UDPconn 中
			_, err := UDPconn.Write(data)
			if err != nil {
				fmt.Println(err)
				return

			}
		}
	}

}

// 完成数据接收协程
func udpRecvProc() {
	fmt.Println("udpRecvProc...")
	UDPconn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(UDPconn)
	if err != nil {
		fmt.Println(err)
	}
	for {
		var buf [512]byte
		// 非定向读取数据到 data 中
		n, err := UDPconn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		// 消息分发
		disPatch(buf[0:n])

	}

}

// 根据 targetID,消息类型进行消息分发
func disPatch(data []byte) {
	fmt.Println("disPatch begins...")
	// 初始化 message
	// 初始化,需要对数据的接受者进行绑定
	msg := Message{}
	// 反序列化到 msg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch msg.Type {
	case 1: // 用户对用户
		sendMsg(msg.TargetId, data)
	case 2:
		sendGroupMsg(msg.TargetId, data)
		//case 3:
		//	// 更新心跳
		//	node := clientMap[msg.UserId]
		//	currentTime := uint64(time.Now().Unix())
		//	node.Heartbeat(currentTime)

	}
	fmt.Println("disPatch ends...")
}
func sendMsg(targetId string, data []byte) {
	fmt.Println("sendMsg...")
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()

	jsonMsg := Message{}
	err := json.Unmarshal(data, &jsonMsg)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	// 将 userId 转成 String
	targetIdStr := strconv.Itoa(int(targetId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	// 刷新消息创建时间
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	// 从Redis 中获取在线信息
	r, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Println(err)
	}

	if r != "" {
		// 如果获取到了,说明用户在线
		if ok {
			fmt.Println("sendMsg >>> userID: ", targetId, "  msg:", string(data))
			node.DataQueue <- data
		}
	}
	// 对 key 做一些处理
	var key string
	if targetId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1
	ress, e := utils.Red.ZAdd(ctx, key, &redis.Z{score, data}).Result() //jsonMsg
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //备用 后续拓展 记录完整msg
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)

}

// 需要重写此方法才能完整的msg转byte[]

func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// 获取缓存里面的消息

func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	rwLocker.RLock()
	//node, ok := clientMap[userIdA]
	rwLocker.RUnlock()
	//jsonMsg := Message{}
	//json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	//key = "msg_" + userIdStr + "_" + targetIdStr
	//rels, err := utils.Red.ZRevRange(ctx, key, 0, 10).Result()  //根据score倒叙

	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	// 发送推送消息
	/**
	// 后台通过websoket 推送消息
	for _, val := range rels {
		fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
		node.DataQueue <- []byte(val)
	}**/
	return rels
}

// 群发就是给每一个人都发消息
func sendGroupMsg(targetId string, msg []byte) {
	userIds := SearchUsersByGroupId(targetId)
	for i := 0; i < len(userIds); i++ {
		//排除给自己的
		if targetId != userIds[i].UID {
			sendMsg(userIds[i].UID, msg)
		}

	}
}

// 更新用户心跳

func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}

// 清理超时连接

func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()

	currentTime := uint64(time.Now().Unix())
	for i := range clientMap {
		node := clientMap[i]
		if node.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("心跳超时..... 关闭连接：", node)
			node.Conn.Close()
		}
	}
	return result
}

// 用户心跳是否超时

func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Println("您因长时间不发消息,为确保安全,已自动下线", node)
		timeout = true
	}
	return
}
