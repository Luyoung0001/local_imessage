package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"html/template"
	"local_imessage/models"
	"local_imessage/utils"
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
// @Summary 用户登录
// @Tags 用户模块
// @Param phone formData string false "手机号码"
// @Param password formData string false "密码"
// @Success 200 {string} json{"code": 0, "message": "登录成功", "data": UserBasic}
// @Failure 200 {string} json{"code": -1, "message": "登录失败"}
// @Router /user/login [post]
func Login(c *gin.Context) {
	// 拿到前端传来的用户名和密码
	phone := c.Request.FormValue("phone")
	password := c.Request.FormValue("password")
	// 查询改用户是否存在
	user := models.FindUserByPhone(phone)
	// 如果不存在
	if user.Phone == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "该用户不存在",
		})
	} else {
		flag := utils.ValidPassword(password, user.Salt, user.PassWord)
		if !flag {
			c.JSON(200, gin.H{
				"code":    -1, //  0成功   -1失败
				"message": "密码不正确",
			})
		} else {
			ind, _ := template.ParseFiles(viper.GetString("path.login"))
			// 传递 userId 到模板
			// 渲染模板并传递参数
			_ = ind.ExecuteTemplate(c.Writer, viper.GetString("path.login"), user.UID)

			// 传入uid
			_ = ind.Execute(c.Writer, user.UID)

		}

	}
}

// UnRegister
// @Summary 注销账号
// @Tags 用户模块
// @param userId formData string false "userId"
// @param password formData string false "password"
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
// @Success 200 {string} welcome
// @Router /chat [get]
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
