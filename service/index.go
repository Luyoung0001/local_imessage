package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"html/template"
)

// GetIndex
// @Summary 获取首页
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles(viper.GetString("path.index"))
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, viper.GetString("notice.news"))
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
	ind, err := template.ParseFiles(viper.GetString("path.toRegister"))
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
	// 这个页面将会非常复杂,这个页面是用户登陆成功后的跳转页面,里面集成了聊天,群聊,好友列表,在线列表,个人信息设置,群管理等等
	// 建议在这个页面下建立一个路由表
	ind, err := template.ParseFiles(viper.GetString("path.login"))
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
	ind, err := template.ParseFiles(viper.GetString("path.unRegister"))
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "unRegister")
	if err != nil {
		return
	}

}

// Chat
// @Tags 用户模块
// @param userId1 formData string false "userId1"
// @param userId2 formData string false "userId2"
// @Success 200 {string} welcome
// @Router /chat [get]
func Chat(c *gin.Context) {

}
