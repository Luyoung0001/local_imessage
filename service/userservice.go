package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"local_imessage/models"
	"local_imessage/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	// 从数据库获得数据,将所有的数据存储成数据,然后返回
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name formData string false "name"
// @param password formData string false "password"
// @param Identity formData string false "Identity"
// @param phone formData string false "phone"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 先判断是否有冲突
	user.Name = c.Request.FormValue("name")
	user.OldPhone = c.Request.FormValue("phone")

	if !models.IsUnique(user) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "请检查你的用户名或者电话,邮箱,它们已被注册!",
		})
		return
	}
	// 判断是否输入完了账号,密码
	passWord := c.Request.FormValue("password")
	rePassWord := c.Request.FormValue("Identity")
	if user.Name == "" || passWord == "" || rePassWord == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码或者账号不能为空!",
		})
		return
	}
	// 如果一切正常,继续
	// 获得一个随机数
	salt := fmt.Sprintf("%06d", rand.Int31())
	if passWord != rePassWord {
		c.JSON(-4, gin.H{
			"code":    -1,
			"message": "两次密码不一致!",
		})

	} else {
		user.Salt = salt
		// 这里暂时存入一个不准确的时间

		user.PassWord = utils.MakePassword(passWord, user.Salt)
		user.HeartBeatTime = time.Now()
		user.LoginTime = time.Now()
		user.LoginOutTime = time.Now()

		re := models.CreateUser(user)
		if re == true {
			c.JSON(200, gin.H{
				"code":    0,
				"message": "新增用户成功!",
				"data":    user,
			})

		} else {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": "新增用户失败!",
			})
		}

	}

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param userId formData string false "userId"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找\删除操作
	uid := c.Query("userId")
	user.UID = uid
	re := models.DeleteUser(user)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "删除用户成功!",
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "删除用户失败!",
		})

	}

}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param userId formData string false "userId"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找,删除操作
	uid := c.PostForm("UID")
	user.UID = uid

	user.Name = c.PostForm("name")
	user.NewPhone = c.PostForm("phone")

	// 判断相异性
	if !models.IsUnique(user) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "请检查你的手机号码",
		})
		return
	}
	// 这个要和 salt 一起更新,这个很关键,不然修改密码后,无法登陆
	// 获得一个随机数
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	passwordRaw := c.PostForm("password")
	PWD := utils.MakePassword(passwordRaw, user.Salt)
	user.PassWord = PWD

	re := models.UpdateUser(user)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0, // 0 : 成功; -1 : 失败
			"message": "修改用户成功!",
			"data":    user,
		})

	} else {
		c.JSON(200, gin.H{
			"code":    -1, // 0 : 成功; -1 : 失败
			"message": "修改用户失败!",
		})
	}

}

// FindUserByPhoneAndPwd
// @Summary 用户登陆
// @Tags 用户模块
// @param phone formData string false "phone"
// @param password formData string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByPhoneAndPwd [post]
func FindUserByPhoneAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	// 拿到前端传来的用户名和密码

	phone := c.Request.FormValue("phone")
	password := c.Request.FormValue("password")
	// 查询改用户是否存在
	user := models.FindUserByPhone(phone)
	// 如果不存在
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "该用户不存在",
			"data":    data,
		})
	} else {
		flag := utils.ValidPassword(password, user.Salt, user.PassWord)
		if !flag {
			c.JSON(200, gin.H{
				"code":    -1, //  0成功   -1失败
				"message": "密码不正确",
				"data":    data,
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code":    0, //  0成功   -1失败
				"message": "登录成功",
				"data":    data,
			})
		}

	}

}

// 防止跨站域的伪造请求

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SendMsg
// @Summary 发送消息
// @Tags 消息模块
// @Success 200 {string} json{"code","message"}
// @Router /user/sendMsg [get]
func SendMsg(c *gin.Context) {
	// 普通的 HTTP 连接升级为 WebSocket 连接
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)

	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("hello, 发送消息呢:", msg)
	// 选取时间格式
	tm := time.Now().Format("2006-01-02 15:15:03")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)

	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}

}

// SendUserMsg
// @Summary 发送消息
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/sendUserMsg [get]
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)

}

// FriendList
// @Summary 好友列表
// @Tags 用户模块
// @param username formData string false "phone"
// @Success 200 {string} json{"code","message"}
// @Router /searchFriends [post]
func FriendList(c *gin.Context) {

}

// AddFriend
// @Summary 添加好友
// @Tags 用户模块
// @param userId formData string false "userId"
// @param targetId formData string false "targetId"
// @Success 200 {string} json{"code","message"}
// @Router /contact/addFriend [post]
func AddFriend(c *gin.Context) {

}

// CreateGroup
// @Summary 创建群聊
// @Tags 用户模块
// @param ownerId formData string false "ownerId"
// @param name formData string false "name"
// @param icon formData string false "icon"
// @param desc formData string false "desc"
// @Success 200 {string} json{"code","message"}
// @Router /contact/createGroup [post]
func CreateGroup(c *gin.Context) {

}

// LoadGroup
// @Summary 加载群聊列表
// @Tags 用户模块
// @param ownerId formData string false "ownerId"
// @Success 200 {string} json{"code","message"}
// @Router /contact/loadcommunity [post]
func LoadGroup(c *gin.Context) {
}

// JoinGroup
// @Summary 加入群聊
// @Tags 用户模块
// @param userId formData string false "userId"
// @param comId formData string false "comId"
// @Success 200 {string} json{"code","message"}
// @Router /contact/joinGroup [post]
func JoinGroup(c *gin.Context) {

}

// RedisMsg
// @Summary redis收发消息
// @Tags 用户模块
// @param userIdA formData string false "userIdA"
// @param userIdB formData string false "userIdB"
// @param start formData string false "start"
// @param end formData string false "end"
// @param isRev formData string false "isRev"
// @Success 200 {string} json{"code","message"}
// @Router /user/find [post]
func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

func UnRegister(c *gin.Context) {

}
func FindPassword(c *gin.Context) {

}
func DeleteFriend(c *gin.Context) {

}
func FriendsStatus(c *gin.Context) {

}
func BlockFriend(c *gin.Context) {

}
func DeleteGroup(c *gin.Context) {

}
func OutGroup(c *gin.Context) {

}
func OwnerManGroup(c *gin.Context) {

}
func ManManGroup(c *gin.Context) {

}
func AllowGroup(c *gin.Context) {

}
