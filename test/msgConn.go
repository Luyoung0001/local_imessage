package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"imessage/models"
)

// 这里就是将模板导入数据库的一个简单案例
// 前提是创建好了一个数据库,数据库的名字暂定义为 "msg"
// 这个是用来实现模型绑定到数据的功能,只运行一次

var (
	DBMsg *gorm.DB
)

func initMySQLMsg() (err error) {
	// 连接数据库
	dsn := "root:791975457@qq.com@tcp(127.0.0.1:3306)/msg?charset=utf8mb4&parseTime=True&loc=Local"
	DBMsg, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 判断是否连通
	return DBMsg.DB().Ping()
}
func main() {
	// 连接数据库
	err := initMySQLMsg()
	if err != nil {
		panic(err)
	}
	// 模型绑定
	DBMsg.AutoMigrate(&models.Message{}) // todos
	defer func(DB *gorm.DB) {
		err := DB.Close()
		if err != nil {
			panic(err)
		}
	}(DBMsg)

}
