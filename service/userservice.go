package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"local_imessage/models"
	"local_imessage/utils"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type ResponseObject struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// GetUserList
// @Summary 获取所有用户
// @Tags 用户模块
// @Success 200 {string} ResponseObject
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	// 从数据库获得数据,将所有的数据存储成数据,然后返回
	data := make([]models.UserBasic, 0)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param username formData string false "用户名"
// @Param password formData string false "密码"
// @Param Identity formData string false "确认密码"
// @Param phone formData string false "电话号码"
// @Success 200 {string} json{"code": 0, "message": "新增用户成功!", "data": UserBasic}
// @Failure 200 {string} json{"code": -1, "message": "新增用户失败!"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 先判断是否有冲突
	user.Name = c.Request.FormValue("username")
	user.Phone = c.Request.FormValue("phone")
	fmt.Println(user.Name, user.Phone)

	if !models.IsUnique(user) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "请检查你的电话,它们已被注册!",
		})
		return
	}
	// 判断是否输入完了账号,密码
	passWord := c.Request.FormValue("password")
	rePassWord := c.Request.FormValue("Identity")

	if user.Name == "" || user.Phone == "" || passWord == "" || rePassWord == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码账号以及用户名不能为空!",
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
// @Param userId formData string false "用户ID"
// @Success 200 {string} json{"code": 0, "message": "删除用户成功!"}
// @Failure 200 {string} json{"code": -1, "message": "删除用户失败!"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找\删除操作
	uid := c.Request.FormValue("userId")
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
// @Summary 修改昵称和密码
// @Tags 用户模块
// @Param userId formData string false "用户ID"
// @Param name formData string false "用户名"
// @Param password formData string false "密码"
// @Success 200 {string} json{"code": 0, "message": "修改用户成功!", "data": UserBasic}
// @Failure 200 {string} json{"code": -1, "message": "修改用户失败!"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找,删除操作
	uid := c.Request.FormValue("userId")
	user.Name = c.Request.FormValue("name")
	user.UID = uid

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

// FriendsList
// @Summary 好友列表
// @Tags 用户模块
// @Param userId formData string false "用户ID"
// @Success 200 {string} json{"code": 0, "message": "获取列表成功!", "data": []UserBasic}
// @Failure 200 {string} json{"code": -1, "message": "获取列表失败"}
// @Router /friendsList [post]
func FriendsList(c *gin.Context) {
	uid := c.Request.FormValue("userId")
	// 通过 uid 获取到用户,查询该用户的 contact 信息
	// 从 contact 信息,映射到 userBasic
	currentUser := models.FindUserByUID(uid)
	data := models.FriendsList(currentUser)
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "获取列表成功!",
		"data":    data,
	})
}

func FindPassword(c *gin.Context) {

}

// AddFriend
// @Summary 添加好友
// @Tags 用户模块
// @Param userId formData string false "发起请求的用户ID"
// @Param targetId formData string false "目标用户ID"
// @Success 200 {string} json{"code": 0, "message": "添加成功"}
// @Failure 200 {string} json{"code": -1, "message": "添加失败"}
// @Router /contact/addFriend [post]
func AddFriend(c *gin.Context) {
	userId := c.Request.FormValue("userId")
	targetId := c.Request.FormValue("targetId")
	re := models.AddFriend(userId, targetId)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0, //  0成功   -1失败
			"message": "添加成功",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "添加失败!",
		})
	}

}

// CreateGroup
// @Summary 创建群聊
// @Tags 用户模块
// @Param ownerId formData string false "群主用户ID"
// @Param groupName formData string false "群组名称"
// @Success 200 {string} json{"code": 0, "message": "创建成功!"}
// @Failure 200 {string} json{"code": -1, "message": "创建失败!"}
// @Router /contact/createGroup [post]
func CreateGroup(c *gin.Context) {
	var group models.GroupBasic

	group.OwnerUID = c.Request.FormValue("ownerId")
	group.Name = c.Request.FormValue("groupName")

	re := models.CreatGroup(group)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0, //  0成功   -1失败
			"message": "创建成功!",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "创建失败!",
		})
	}
}

// GroupsList
// @Summary 加载群聊列表
// @Tags 用户模块
// @Param userId formData string false "用户ID"
// @Success 200 {string} json{"code": 0, "message": "获取成功!", "data": []GroupBasic}
// @Failure 200 {string} json{"code": -1, "message": "获取失败"}
// @Router /groupsList [post]
func GroupsList(c *gin.Context) {
	uid := c.Request.FormValue("userId")
	groupList := models.GroupsList(uid)
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "获取成功!",
		"data":    groupList,
	})

}

// JoinGroup
// @Summary 加入群聊
// @Tags 用户模块
// @Param userId formData string false "用户ID"
// @Param groupId formData string false "群组ID"
// @Success 200 {string} json{"code": 0, "message": "加入成功!"}
// @Failure 200 {string} json{"code": -1, "message": "加入失败!"}
// @Router /contact/joinGroup [post]
func JoinGroup(c *gin.Context) {
	userId := c.Request.FormValue("userId")
	groupId := c.Request.FormValue("groupId")
	re := models.JoinGroup(userId, groupId)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0, //  0成功   -1失败
			"message": "加入成功!",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "添加失败!",
		})

	}
}

// DeleteFriend
// @Summary 删除好友
// @Tags 用户模块
// @Param userId1 formData string false "发起请求的用户ID"
// @Param userId2 formData string false "待删除的好友ID"
// @Success 200 {string} json{"code": 0, "message": "删除成功!"}
// @Failure 200 {string} json{"code": -1, "message": "删除失败!"}
// @Router /contact/deleteFriend [post]
func DeleteFriend(c *gin.Context) {
	userId1 := c.Request.FormValue("userId1")
	userId2 := c.Request.FormValue("userId2")
	// 直接删除关系
	re := models.DeleteFriend(userId1, userId2)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0, //  0成功   -1失败
			"message": "删除成功!",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "删除失败!",
		})
	}

}

// FriendsOnlineList
// @Summary 好友在线列表
// @Tags 用户模块
// @param userId formData string false "userId"
// @Success 200 {string} json{"code","message"}
// @Router /contact/friendsOnlineList [post]
func FriendsOnlineList(c *gin.Context) {
	// 怎么判断在线?
	// 用心跳包

}

// BlockFriend
// @Summary 拉黑好友
// @Tags 用户模块
// @Param userId1 formData string false "发起请求的用户ID"
// @Param userId2 formData string false "待删除的好友ID"
// @Success 200 {string} json{"code": 0, "message": "删除成功!"}
// @Failure 200 {string} json{"code": -1, "message": "删除失败!"}
// @Router /contact/blockFriend [post]
func BlockFriend(c *gin.Context) {
	// 拉黑的本质是什么?
	// 定义拉黑为删除
	DeleteFriend(c)
}

// DeleteGroup
// @Summary 群主解散群
// @Tags 用户模块
// @Param operator formData string true "操作者的ID"
// @Param groupId formData string true "群组ID"
// @Success 200 {string} ResponseObject
// @Router /contact/deleteGroup [post]
func DeleteGroup(c *gin.Context) {
	operator := c.Request.FormValue("operator")
	groupId := c.Request.FormValue("groupId")

	// 执行删除群的操作
	re := models.DeleteGroup(operator, groupId)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0, //  0成功   -1失败
			"message": "删除成功!",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "删除失败!",
		})
	}

}

// DeMemberFromGroup
// @Summary 从群聊中删除成员
// @Tags 用户模块
// @Param operator formData string false "操作者的ID"
// @Param groupId formData string false "群组ID"
// @Param userId formData string false "待删除成员的ID"
// @Success 200 {string} json{"code","message"}
// @Router /contact/deMemberFromGroup [post]
func DeMemberFromGroup(c *gin.Context) {
	operator := c.Request.FormValue("operator")
	groupId := c.Request.FormValue("groupId")
	userId := c.Request.FormValue("userId")
	// operator 在 groupId 中删除 userId
	re := models.LevelUserInGroup(operator, groupId, userId)
	if re == true {
		c.JSON(200, gin.H{
			"code":    0, //  0成功   -1失败
			"message": "删除成功!",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "删除失败!",
		})
	}

}

// AddMan
// @Summary 群主添加userId为群管理员
// @Tags 用户模块
// @Param operator formData string false "操作者的ID"
// @Param groupId formData string false "群组ID"
// @Param userId formData string false "待添加管理员的ID"
// @Success 200 {string} json{"code","message"}
// @Router /contact/addMan [post]
func AddMan(c *gin.Context) {
	// 只有群主才可以添加管理员
	operator := c.Request.FormValue("operator")
	groupId := c.Request.FormValue("groupId")
	userId := c.Request.FormValue("userId")
	relation := models.RelationBetweenUserAndGroup(operator, groupId)
	if relation == 2 {
		// 群主才可以操作
		re := models.AddMan(userId, groupId)
		if re == true {
			c.JSON(200, gin.H{
				"code":    0, //  0成功   -1失败
				"message": "添加成功!",
			})
		} else {
			c.JSON(200, gin.H{
				"code":    -1, //  0成功   -1失败
				"message": "添加失败!",
			})
		}
	} else {
		c.JSON(200, gin.H{
			"code":    -2, //  0成功   -1失败
			"message": "你无权操作!",
		})
	}

}

// OutGroup
// @Summary 退群
// @Tags 用户模块
// @Param userId formData string false "用户ID"
// @Param groupId formData string false "群组ID"
// @Success 200 {string} json{"code","message"}
// @Router /contact/outGroup [post]
func OutGroup(c *gin.Context) {
	groupId := c.Request.FormValue("groupId")
	userId := c.Request.FormValue("userId")
	// 任何人都有资格退群,群处除外
	relation := models.RelationBetweenUserAndGroup(userId, groupId)
	// 只有成员或者管理员才可以退群
	if relation == 1 || relation == 0 {
		//
		re := models.OutGroup(userId, groupId)
		if re == true {
			c.JSON(200, gin.H{
				"code":    0, //  0成功   -1失败
				"message": "退群成功!",
			})
		} else {
			c.JSON(200, gin.H{
				"code":    -1, //  0成功   -1失败
				"message": "退群失败!",
			})
		}

	} else {
		c.JSON(200, gin.H{
			"code":    -2, //  0成功   -1失败
			"message": "群主不能退群!",
		})
	}

}

// GetContactList
// @Summary 获取所有关系
// @Tags 用户模块
// @Success 200 {string} ResponseObject
// @Router /contact/getContactList [get]
func GetContactList(c *gin.Context) {
	// 从数据库获得数据,将所有的数据存储成数据,然后返回
	data := make([]models.Contact, 0)
	data = models.GetContactList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": data,
	})

}

// GetGroupList
// @Summary 获取群的列表
// @Tags 用户模块
// @Success 200 {string} ResponseObject
// @Router /contact/getGroupList [get]
func GetGroupList(c *gin.Context) {
	data := make([]models.GroupBasic, 0)
	data = models.GetGroupList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": data,
	})
}

// SearchUsersByGroupId
// @Summary 按照groupId 找到所有的群成员
// @Tags 用户模块
// @Param groupId formData string false "群组ID"
// @Success 200 {string} json{"code","message"}
// @Router /contact/searchUsersByGroupId [post]
func SearchUsersByGroupId(c *gin.Context) {
	groupId := c.Request.FormValue("groupId")
	data := models.SearchUsersByGroupId(groupId)
	fmt.Println(data)
	c.JSON(200, gin.H{
		"code":    0,
		"message": data,
	})

}

// ListInOnePage
// @Summary 给前端呈现列表
// @Tags 界面模块
// @Param userId formData string false "用户ID"
// @Success 200 {string} json{"owner": {...}, "contactList": [...], "onLineList": [...], "groupList": [...]}
// @Router /listInOnePage [post]
func ListInOnePage(c *gin.Context) {
	userId := c.PostForm("userId") // 使用 PostForm 获取表单值
	// 打印一下
	fmt.Println(userId)
	data, err := models.ListInOnePage(userId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching user info"})
		return
	}
	fmt.Println(data.Owner)
	fmt.Println(data.GroupList)
	fmt.Println(data.ContactList)
	fmt.Println(data.OnLineList)
	c.JSON(http.StatusOK, data)
}
