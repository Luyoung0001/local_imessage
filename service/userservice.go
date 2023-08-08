package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"imessage/models"
	"imessage/utils"
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
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}

	// 先判断是否有冲突
	user.Name = c.Request.FormValue("name")
	user.Phone = c.Request.FormValue("phone")
	user.Email = c.Request.FormValue("email")

	if !models.IsUniqueCreateUser(user) {
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

		models.CreateUser(user)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "新增用户成功!",
			"data":    user,
		})
	}

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找\删除操作
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除用户成功!",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找,删除操作
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)

	user.Name = c.PostForm("name")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	// 判断相异性

	if !models.IsUniqueUpdateUser(user) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "请检查你的用户名或者电话,邮箱,它们已被注册!",
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

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, // 0 : 成功; -1 : 失败
			"message": "修改用户失败!",
		})

	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, // 0 : 成功; -1 : 失败
			"message": "修改用户成功!",
			"data":    user,
		})
	}

}

// FindUserByNameAndPwd
// @Summary 用户登录
// @Tags 用户模块
// @param name formData string false "name"
// @param password formData string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	// 拿到前端传来的用户名和密码

	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	fmt.Println(name, password)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "登录成功",
		"data":    data,
	})
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

// SearchFriends
// @Summary 搜索好友
// @Tags 用户模块
// @param userId formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /searchFriends [post]
func SearchFriends(c *gin.Context) {

	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))
	// 这里给前端返回一个请求头,里面包含好友列表 users
	utils.RespOKList(c.Writer, users, len(users))
}

// AddFriend
// @Summary 添加好友
// @Tags 用户模块
// @param userId formData string false "userId"
// @param targetName formData string false "targetName"
// @Success 200 {string} json{"code","message"}
// @Router /contact/addFriend [post]
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// CreateCommunity
// @Summary 创建群聊
// @Tags 用户模块
// @param ownerId formData string false "ownerId"
// @param name formData string false "name"
// @param icon formData string false "icon"
// @param desc formData string false "desc"
// @Success 200 {string} json{"code","message"}
// @Router /contact/createCommunity [post]
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	// 创建好了群,之后要重新刷新群的列表
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// LoadCommunity
// @Summary 加载好友列表
// @Tags 用户模块
// @param ownerId formData string false "ownerId"
// @Success 200 {string} json{"code","message"}
// @Router /contact/loadcommunity [post]
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	// 把 data 查找出来发给客户端
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// JoinGroups
// @Summary 加入群聊
// @Tags 用户模块
// @param userId formData string false "userId"
// @param comId formData string false "comId"
// @Success 200 {string} json{"code","message"}
// @Router /contact/joinGroup [post]
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// FindByID
// @Summary 查看用户信息
// @Tags 用户模块
// @param userId formData string false "userId"
// @Success 200 {string} json{"code","message"}
// @Router /user/find [post]
func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
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
func DeleteCommunity(c *gin.Context) {

}
func OutCommunity(c *gin.Context) {

}
func OwnerManCommunity(c *gin.Context) {

}
func ManManCommunity(c *gin.Context) {

}
func AllowCommunity(c *gin.Context) {

}
