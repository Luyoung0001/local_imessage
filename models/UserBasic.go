package models

import (
	"context"
	"encoding/json"
	"fmt"
	"local_imessage/utils"
	"log"
	"time"
)

// 设计原则:底层干最简单的事情,将控制全部放在上层,上层拥有最高的灵活度和可设计性

type UserBasic struct {
	// 加密后的字符串
	UID           string // 由Phone 生成,且永远不变,用户唯一表示符
	Phone         string
	Name          string
	PassWord      string
	Identity      string
	ClientIP      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartBeatTime time.Time
	LoginOutTime  time.Time
	IsLogOut      bool
	DeviceInfo    string
	DataType      string `json:"dataType"` // 添加 DataType 字段来标识数据类型
}

// 返回 name

func (table *UserBasic) TableName() string {
	return "userBasic"
}

// 返回所有的用户列表
// 为各种函数服务

func GetUserList() []UserBasic {
	userLists := make([]UserBasic, 0)
	ctx := context.Background()
	var keys []string
	// 获取所有键
	keys, err := utils.Red.Keys(ctx, "*").Result()
	if err != nil {
		log.Fatal(err)
	}
	// 遍历每个键，获取值
	for _, key := range keys {
		value, err := utils.Red.Get(ctx, key).Result()
		if StructType(value) == "userBasic" {
			var userInfo UserBasic
			if err != nil {
				log.Println("Error:", err)
				continue
			}
			// 反序列化
			err = json.Unmarshal([]byte(value), &userInfo)
			if err != nil {
				return nil
			}
			// 添加到 values
			userLists = append(userLists, userInfo)
		}

	}
	return userLists
}

// 通过用户姓名定位到一群人

func FindUserByName(name string) []UserBasic {
	userLists := GetUserList()
	var nameLists []UserBasic
	for _, user := range userLists {
		if user.Name == name {
			nameLists = append(nameLists, user)
		}
	}
	return nameLists
}

// 通过电话号码定位到一个人

func FindUserByPhone(phone string) UserBasic {
	// 为了保护隐私 只能通过新电话定位到一个人
	userLists := GetUserList()
	for _, user := range userLists {
		if user.Phone == phone {
			return user
		}
	}
	return UserBasic{}
}

// 创建用户

func CreateUser(user UserBasic) bool {
	// 这个用户只传回了
	// user.phone
	// user.Name
	// user.PassWord

	ctx := context.Background()
	user.UID = IdGenerator()
	user.DataType = "userBasic"

	// 将用户对象序列化为 JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 存储用户数据到 Redis
	err = utils.Red.Set(ctx, user.UID, userJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func DeleteUser(user UserBasic) bool {
	userKey := user.UID
	return DeleteKvFromRed(userKey)
}

// 更新用户信息

func UpdateUser(user UserBasic) bool {
	ctx := context.Background()
	// 序列化更新后的结构体为 JSON
	updatedJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 存储更新后的值回 Redis
	err = utils.Red.Set(ctx, user.UID, updatedJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// 只在创建时进行校验
// 用户更新自己的信息时,暂时不可更新电话号码

func IsUnique(user UserBasic) bool {
	userLists := GetUserList()
	for _, item := range userLists {
		if user.Phone == item.Phone {
			return false
		}
	}
	return true

}

// 按照 UID 进行查询

func FindUserByUID(uid string) UserBasic {
	userLists := GetUserList()
	for _, user := range userLists {
		if user.UID == uid {
			return user
		}
	}
	return UserBasic{}
}

// 传入一个 User,返回该用户的好友列表[]UserBasic

func FriendsList(user UserBasic) []UserBasic {
	// 拿到 contacts
	friendsList := make([]UserBasic, 0)
	contactList := FriendList(user.UID)
	fmt.Println(contactList)
	// 获取
	for _, each := range contactList {
		friend := FindUserByUID(each.TargetId)
		friendsList = append(friendsList, friend)
	}
	return friendsList
}
