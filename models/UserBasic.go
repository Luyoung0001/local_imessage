package models

import (
	"fmt"
	"gorm.io/gorm"
	"imessage/utils"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIP      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartBeatTime time.Time
	LoginOutTime  time.Time
	IsLogOut      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)

	//for _, v := range data {
	//	fmt.Println(v)
	//}
	return data
}

// 通过密码和用户名来查询用户

func FindUserByNameAndPwd(name, password string) UserBasic {
	// password 是加密后的字符串
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user)

	// token 加密

	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)

	utils.DB.Where("name = ?", user.Name).Update("identity", temp)
	utils.DB.Model(&user).Update("identity", temp)
	return user
}

// 通过用户姓名定位到一个人

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// 通过电话号码定位到一个人

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

// 通过邮箱定位到一个人

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email, Salt: user.Salt})
}
func IsUniqueUpdateUser(user UserBasic) bool {
	data1 := FindUserByName(user.Name)
	data2 := FindUserByPhone(user.Phone)
	data3 := FindUserByEmail(user.Email)
	// 同时还得检查被搜到的 id 是否和自己一样,如果一样,那就说明可以修改;否则就不能修改
	if data1.Name != "" && data1.ID != user.ID || data2.Phone != "" && data2.ID != user.ID || data3.Email != "" && data3.ID != user.ID {
		return false
	}
	return true
}
func IsUniqueCreateUser(user UserBasic) bool {
	data1 := FindUserByName(user.Name)
	data2 := FindUserByPhone(user.Phone)
	data3 := FindUserByEmail(user.Email)
	if data1.Name != "" || data2.Phone != "" || data3.Email != "" {
		return false
	}
	return true
}

func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}
