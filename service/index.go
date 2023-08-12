package service

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"local_imessage/models"
)

// GetIndex
// @Summary 获取首页
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("/Users/luliang/GoLand/local_imessage/index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		return
	}
}

// ToRegister
// @Summary 跳转到注册页面
// @Tags 首页
// @Success 200 {string} welcome
// @Router /toRegister [post]
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("/Users/luliang/GoLand/local_imessage/views/user/register.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "register")
	if err != nil {
		return
	}

}

// Login
// @Tags 登录
// @Success 200 {string} welcome
// @Router /login [get]
func Login(c *gin.Context) {
	ind, err := template.ParseFiles("/Users/luliang/GoLand/local_imessage/views/user/register.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "login")
	if err != nil {
		return
	}

}

// UnRegister
// @Summary 注销账号
// @Tags 用户模块
// @param userId formData string false "userId"
// @param groupId formData string false "groupId"
// @Success 200 {string} json{"code","message"}
// @Router /contact/unRegister [post]
func UnRegister(c *gin.Context) {
	ind, err := template.ParseFiles("/Users/luliang/GoLand/local_imessage/views/user/unregister.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "login")
	if err != nil {
		return
	}

}

// ToChat
// @Tags 用户模块
// @param UserId query string false "userid"
// @param token query string false "token"
// @Success 200 {string} welcome
// @Router /toChat [get]
func ToChat(c *gin.Context) {

}

// Chat
// @Tags 用户模块
// @Success 200 {string} welcome
// @Router /chat [get]
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
