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

var userLists []*UserBasic

type UserBasic struct {
	// 加密后的字符串
	UID           string // 由OldPhone 生成,且永远不变,用户唯一表示符
	OldPhone      string
	NewPhone      string
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
}

// 返回 name

func (table *UserBasic) TableName() string {
	return "user_basic"
}

// 返回所有的用户列表

func GetUserList() []*UserBasic {
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
		userLists = append(userLists, &userInfo)
	}
	return userLists
}

// 通过密码和用户名来查询
// 就是在做登录校验

func FindUserByUIDAndPwd(UID, password string) UserBasic {
	// password 是加密后的字符串
	ctx := context.Background()
	userString, err := utils.Red.Get(ctx, UID).Result()
	if err != nil {
		fmt.Println(err)
		return UserBasic{}
	}
	// 反序列化
	var user UserBasic
	err = json.Unmarshal([]byte(userString), &user)
	if err != nil {
		return UserBasic{}
	}
	// 判断
	if user.PassWord == password {
		return user
	} else {
		fmt.Println("密码错误!")
	}
	return UserBasic{}
}

// 通过用户姓名定位到一群人

func FindUserByName(name string) []*UserBasic {
	userLists = GetUserList()
	var nameLists []*UserBasic
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
	userLists = GetUserList()
	for _, user := range userLists {
		if user.NewPhone == phone {
			return *user
		}
	}
	return UserBasic{}
}

// 创建用户

func CreateUser(user UserBasic) bool {
	// 这个用户只传回了
	// user.OldPhone
	// user.Name
	// user.PassWord
	user.NewPhone = user.OldPhone

	ctx := context.Background()

	// 获取Phone
	phone := user.OldPhone
	// 创造 userId
	UID := utils.MD5Encode(phone)
	// 将用户对象序列化为 JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 存储用户数据到 Redis
	err = utils.Red.Set(ctx, UID, userJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func DeleteUser(user UserBasic) bool {
	ctx := context.Background()

	// 要删除的用户键
	userKey := user.UID

	// 删除用户
	deleted, err := utils.Red.Del(ctx, userKey).Result()
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	if deleted > 0 {
		fmt.Printf("User with key '%s' deleted.\n", userKey)
	} else {
		fmt.Printf("User with key '%s' not found.\n", userKey)
		return false
	}
	return true
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
func IsUnique(user UserBasic) bool {
	currentPhone := user.NewPhone
	// 要更新一个值,就要判断要更新的新的电话号码是否已经被别的用户使用
	// 因此这里只判断电话号码是否被注册过(查看 NewPhone 即可,OldPhone 已被校验)
	// 遍历
	// 初始化游标
	userLists = GetUserList()
	for _, item := range userLists {
		if item.NewPhone == currentPhone {
			return false
		}
	}
	return true

}

// 按照 UID 进行查询

func FindUserByUID(uid string) UserBasic {
	userLists = GetUserList()
	for _, user := range userLists {
		if user.UID == uid {
			return *user
		}
	}
	return UserBasic{}
}
