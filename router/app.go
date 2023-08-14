package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"local_imessage/docs"
	"local_imessage/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = ""

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 静态资源
	r.Static("/asset", viper.GetString("path.static"))
	r.LoadHTMLGlob(viper.GetString("path.LoadHTMLGlob"))

	// 首页
	r.GET("/", service.GetIndex)                  //网站主页
	r.GET("/index", service.GetIndex)             // 网站主页
	r.POST("/user/login", service.Login)          // 登陆以及登陆后的页面呈现
	r.POST("/toRegister", service.ToRegister)     // 用户注册
	r.POST("/friendsList", service.FriendsList)   // 返回好友列表
	r.POST("/groupsList", service.GroupsList)     // 返回群列表
	r.POST("/findPassword", service.FindPassword) // 利用手机验证码找回密码
	r.POST("/unRegister", service.UnRegister)     // 注销账号
	// 好友管理
	r.POST("/contact/deFriend", service.DeleteFriend)               // 删除好友
	r.POST("/contact/friendsOnlineList", service.FriendsOnlineList) // 实现在线好友
	r.POST("/contact/blockFriend", service.BlockFriend)             // 实现屏蔽好友消息
	r.POST("/contact/addFriend", service.AddFriend)                 // 添加好友页面
	r.GET("/contact/getContactList", service.GetContactList)        // 获取所有关系列表

	// 群
	r.POST("/contact/createGroup", service.CreateGroup)                   // 创建群
	r.POST("/contact/deGroup", service.DeleteGroup)                       // 解散群
	r.POST("/contact/addMan", service.AddMan)                             // 添加管理员
	r.POST("/contact/outGroup", service.OutGroup)                         // 用户退出群聊
	r.POST("/contact/joinGroup", service.JoinGroup)                       // 添加群聊
	r.POST("/contact/deMemberFromGroup", service.DeMemberFromGroup)       // 删除群中的成员
	r.GET("/contact/getGroupList", service.GetGroupList)                  // 获取所有群列表
	r.POST("/contact/searchUsersByGroupId", service.SearchUsersByGroupId) // 按照groupId 找到所有成员

	// 用户模块
	r.POST("/user/createUser", service.CreateUser)                       // 增加用户
	r.GET("/user/getUserList", service.GetUserList)                      // 获取用户列表
	r.POST("/user/deleteUser", service.DeleteUser)                       // 删除用户
	r.POST("/user/updateUser", service.UpdateUser)                       // 更新用户
	r.POST("/user/findUserByPhoneAndPwd", service.FindUserByPhoneAndPwd) // 登陆验证

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
	// 消息模块

	r.POST("chat", service.Chat) // 聊天

	return r
}
