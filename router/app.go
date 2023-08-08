package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"imessage/docs"
	"imessage/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = ""

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 静态资源
	r.Static("/asset", "/Users/luliang/GoLand/imessage/asset/")
	r.LoadHTMLGlob("/Users/luliang/GoLand/imessage/views/**/*")

	// 首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)        // 用户注册
	r.GET("/toChat", service.ToChat)                // 聊天页面
	r.POST("/contact/addFriend", service.AddFriend) // 添加好友页面
	r.POST("/searchFriends", service.SearchFriends) // 返回好友列表
	r.POST("/findPassword", service.FindPassword)   // 利用手机验证码找回密码
	r.POST("/unRegister", service.UnRegister)       // 注销页面
	// 好友管理
	r.POST("/contact/deFriend", service.DeleteFriend)       // 删除好友
	r.POST("/contact/friendsStatus", service.FriendsStatus) // 实现显示好友在线状态
	r.POST("/contact/blockFriend", service.BlockFriend)     // 实现屏蔽好友消息
	// 群
	r.POST("/contact/createCommunity", service.CreateCommunity)     // 创建群
	r.POST("/contact/deCommunity", service.DeleteCommunity)         // 解散群
	r.POST("/contact/outCommunity", service.OutCommunity)           // 用户退出群聊
	r.POST("/contact/ownerManCommunity", service.OwnerManCommunity) // 实现群主对群组管理员的添加和删除
	r.POST("/contact/manManCommunity", service.ManManCommunity)     // 实现群组管理员/群主从群组中移除用户
	r.POST("/contact/allowCommunity", service.AllowCommunity)       // 实现群组管理员,群主批准用户加入群组
	r.POST("/contact/loadcommunity", service.LoadCommunity)         // 群列表
	r.POST("/contact/joinGroup", service.JoinGroups)                // 添加群聊
	// 用户模块
	r.POST("/user/createUser", service.CreateUser)   // 增加用户
	r.POST("/user/getUserList", service.GetUserList) // 获取用户列表
	r.POST("/user/deleteUser", service.DeleteUser)   // 删除用户
	r.POST("/user/updateUser", service.UpdateUser)   // 更新用户
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.POST("/user/login", service.FindUserByNameAndPwd)
	r.POST("/user/find", service.FindByID)

	// 发送消息
	//聊天功能
	//对于 私聊和群组 的聊天功能均需要实现：
	//
	//实现查看历史消息记录
	//实现用户间在线聊天
	//实现在线用户对离线用户发送消息，离线用户上线后获得通知
	//实现文件发送的断点续传（提高）
	//实现在线发送文件
	//实现在线用户对离线用户发送文件，离线用户上线后获得通知/接收（提高）
	//实现用户在线时,消息的实时通知
	//收到好友请求
	//收到私聊
	//收到加群申请
	//...
	r.GET("/user/sendMsg", service.SendMsg) // websocket 测试
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.GET("/chat", service.Chat)
	r.POST("/user/redisMsg", service.RedisMsg)

	return r
}
