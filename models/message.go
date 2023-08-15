package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"local_imessage/utils"
	"net"
	"net/http"
	"sync"
	"time"
)

// 自动调度
func init() {
	go udpSendProc()
	go udpRecvProc()
}

type Message struct {
	MessageId  string // 随机生成
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
	DataType   string `json:"dataType"` // 添加 DataType 字段来标识数据类型
}

type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
	GroupSets     set.Interface
}

//func (table *Message) TableName() string {
//	return "message"
//}

// 一个 userId 绑定一个 *Node
// 引用传递,避免大量的值拷贝

var clientMap = make(map[string]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	userId := query.Get("userId")

	isVALID := true
	// 升级链接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isVALID
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		DataQueue:     make(chan []byte, 500),
		GroupSets:     set.New(set.ThreadSafe),    // 线程安全的
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登录时间
	}

	// 安全绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 两个协程
	go sendProc(node)
	go recvProc(node)

	// 设置用户过期时间
	SetUserOnlineInfo("online_"+userId, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)

}

// 发送消息
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

// 接受消息
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println(err)
		}

		if msg.Type == 3 {

			// 如果是心跳包,更新心跳时间
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {

			// 否则进行信息分发
			disPatch(data)
			broadMsg(data)
		}
	}
}

// 将信息全部存进 udpSendChan
var udpSendChan = make(chan []byte, 1024*32)

func broadMsg(data []byte) {
	udpSendChan <- data

}

// 服务器数据发送协程
func udpSendProc() {
	UDPconn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(10, 30, 0, 159),
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
		select {
		case data := <-udpSendChan:
			_, err := UDPconn.Write(data)
			if err != nil {
				fmt.Println(err)
				return

			}
		}
	}

}

// 服务器数据接收协程
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
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch msg.Type {
	case 1:
		// 用户对用户
		sendMsg(msg.TargetId, data)
	case 2:
		// 群发消息
		sendGroupMsg(msg.TargetId, data)
	case 3:
		// 更新心跳
		node := clientMap[msg.UserId]
		currentTime := uint64(time.Now().Unix())
		node.Heartbeat(currentTime)

	}
}
func sendMsg(targetId string, data []byte) {
	// 拿出 node
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

	targetIdStr := targetId
	userIdStr := jsonMsg.UserId

	// 刷新消息创建时间
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	
	// 判断是否在线
	isOnline := IsOnline(userIdStr)
	if isOnline != false {
		// 如果获取到了,说明用户在线
		if ok {
			// 后台打印消息
			fmt.Println("sendMsg >>> userID: ", targetIdStr, "  msg:", string(data))
			node.DataQueue <- data
		}
	}
	// 生成一个用于作为Redis有序集合的键
	var key string
	if targetIdStr > userIdStr {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}

	// 从有序集合中获取已存储的消息数据
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}

	// 计算一个用于有序集合的成员分数
	score := float64(cap(res)) + 1

	// 添加消息
	ress, err := utils.Red.ZAdd(ctx, key, &redis.Z{score, data}).Result()
	if err != nil {
		fmt.Println(err)
	}
	// 打印消息数量
	fmt.Println(ress)

}

// 需要重写此方法才能完整的msg转byte[]

func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

// 获取缓存里面的消息

func RedisMsg(userIdA string, userIdB string, start int64, end int64, isRev bool) []string {

	ctx := context.Background()
	var key string
	if userIdA > userIdB {
		key = "msg_" + userIdB + "_" + userIdA
	} else {
		key = "msg_" + userIdA + "_" + userIdB
	}

	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err)
	}
	// 发送推送消息

	// 后台通过websoket 推送消息
	//for _, val := range rels {
	//	fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
	//	node.DataQueue <- []byte(val)
	//}
	return rels
}

// 群发就是给每一个人都发消息
func sendGroupMsg(targetId string, msg []byte) {
	// 每个成员都会收到消息
	userBasics := SearchUsersByGroupId(targetId)
	for i := 0; i < len(userBasics); i++ {
		// 不能给自己发
		if targetId != userBasics[i].UID {
			sendMsg(userBasics[i].UID, msg)
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
			fmt.Println("心跳停止,关闭连接：", node)
			err := node.Conn.Close()
			if err != nil {
				return false
			}
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
