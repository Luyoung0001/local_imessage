package service

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"imessage/models"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("/Users/luliang/GoLand/imessage/index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		return
	}
}

// ToRegister
// @Tags 首页
// @Success 200 {string} welcome
// @Router /toRegister [get]
func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("/Users/luliang/GoLand/imessage/views/user/register.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "register")
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
	ind, err := template.ParseFiles("/Users/luliang/GoLand/imessage/views/chat/index.html",
		"/Users/luliang/GoLand/imessage/views/chat/head.html",
		"/Users/luliang/GoLand/imessage/views/chat/foot.html",
		"/Users/luliang/GoLand/imessage/views/chat/tabmenu.html",
		"/Users/luliang/GoLand/imessage/views/chat/concat.html",
		"/Users/luliang/GoLand/imessage/views/chat/group.html",
		"/Users/luliang/GoLand/imessage/views/chat/profile.html",
		"/Users/luliang/GoLand/imessage/views/chat/createcom.html",
		"/Users/luliang/GoLand/imessage/views/chat/userinfo.html",
		"/Users/luliang/GoLand/imessage/views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId")) // int
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId) // 转化成 uint
	user.Identity = token
	err = ind.Execute(c.Writer, user)
	if err != nil {
		return
	}

}

// Chat
// @Tags 用户模块
// @Success 200 {string} welcome
// @Router /chat [get]
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
