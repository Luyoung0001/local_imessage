package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"imessage/models"
)

var (
	DBContact *gorm.DB
)

func initMySQLContact() (err error) {
	// 连接数据库
	dsn := "root:791975457@qq.com@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"
	DBContact, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 判断是否连通
	return DBContact.DB().Ping()
}
func main() {
	// 连接数据库
	err := initMySQLContact()
	if err != nil {
		panic(err)
	}
	// 模型绑定
	DBContact.AutoMigrate(&models.Contact{}) // todos
	defer func(DB *gorm.DB) {
		err := DB.Close()
		if err != nil {
			panic(err)
		}
	}(DBContact)

}
