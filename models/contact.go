package models

import (
	"context"
	"encoding/json"
	"fmt"
	"local_imessage/utils"
	"log"
)

var friendList []*Contact
var contactList []*Contact

type Contact struct {
	OwnerId string // key

	TargetId string // 有人的 ID,有群的 ID
	Type     int    // 1 代表好友关系;2 代表群聊关系
	Desc     string // 描述
}

func (table *Contact) TableName() string {
	return "contact"
}

// 返回所有的列表,这个列表为后面的查询等提供服务
func getContactList() []*Contact {
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
		var contact Contact
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		// 反序列化
		err = json.Unmarshal([]byte(value), &contact)
		if err != nil {
			return nil
		}
		// 添加到 values
		contactList = append(contactList, &contact)
	}
	return contactList

}

// 返回 userId 的所有好友列表

func FriendList(userId string) []*Contact {
	contactList = getContactList()
	for _, contact := range contactList {
		if contact.OwnerId == userId && contact.Type == 1 {
			friendList = append(friendList, contact)
		}
	}
	return friendList
}

// 添加好友的同时,对方也对添加你

func AddFriend(userId string, targetID string) {
	ctx := context.Background()
	// 添加好友的本质就是创建一个contact 记录
	var contact1 Contact
	var contact2 Contact
	// 记录 1
	contact1.Type = 1 // 好友关系
	contact1.OwnerId = userId
	contact1.TargetId = targetID
	//记录 2
	contact2.Type = 1 // 好友关系
	contact2.OwnerId = targetID
	contact2.TargetId = userId

	// 序列化
	userJSON1, err1 := json.Marshal(contact1)
	userJSON2, err1 := json.Marshal(contact2)
	if err1 != nil {
		log.Fatal(err1)
	}
	// 写入 redis 并返回
	err1 = utils.Red.Set(ctx, userId, userJSON1, 0).Err()
	err1 = utils.Red.Set(ctx, targetID, userJSON2, 0).Err()
	if err1 != nil {
		log.Fatal(err1)
	}

}

// 删除好友
// 双向删除,a 删除 b,那么 b 也将自动删除 a

func DeleteFriend(userID, targetID string) {
	//删除好友就是删除一个数据
	//删除的前提是对方已经成为你的好友,这个由上层控制,这里不做校验
	// 因此,这里要删除两次
	ctx := context.Background()
	friendList = FriendList(userID)
	for _, friend := range friendList {
		// 只能删除满足特定条件的记录
		if friend.TargetId == targetID {
			deleted, err := utils.Red.Del(ctx, userID).Result()
			if err != nil {
				log.Println("Error:", err)
				return
			}
			if deleted < 0 {
				fmt.Println("删除失败!")
			}
		}
		break
	}
	// 继续删除
	friendList = FriendList(targetID)
	for _, friend := range friendList {
		if friend.TargetId == userID {
			deleted, err := utils.Red.Del(ctx, targetID).Result()
			if err != nil {
				log.Println("Error:", err)
				return
			}
			if deleted < 0 {
				fmt.Println("删除失败!")
			}
		}
		break
	}

}

// 加群

func JoinGroup(userId string, groupID string) {
	// 加群本质就是将 groupID 添加到contact 中
	// 这里默认只要加群就能成功,群主和管理员可以将该成员删除
	// 创建一个记录 contact,键为userUID,值为 groupID
	// 为了方便拉出群的关系,可以同时 再建立一个 groupID----userID 的映射
	ctx := context.Background()
	// 添加好友的本质就是创建一个contact 记录
	var contact1 Contact
	var contact2 Contact

	contact1.Type = 2 // 群聊关系
	contact2.Type = 2
	contact1.TargetId = groupID
	contact2.TargetId = userId

	// 序列化
	userJSON1, err := json.Marshal(contact1)
	userJSON2, err := json.Marshal(contact2)
	if err != nil {
		log.Fatal(err)
	}
	// 写入 redis 并返回
	err = utils.Red.Set(ctx, userId, userJSON1, 0).Err()
	err = utils.Red.Set(ctx, userId, userJSON2, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

}

// 管理群

func MyGroupList(userId string) {
	// 列出我的所有群

}

var userGroupList []*Contact

// 返回用户 userID 加的所有群

func GroupList(userID string) []*Contact {
	// 就是查找type=2
	contactList = FriendList(userID)
	for _, item := range contactList {
		if item.Type == 2 {
			userGroupList = append(userGroupList, item)
		}
	}
	return userGroupList
}

// 退群

func OutGroup(userId, groupId string) {
	// 双相删除记录
	DeleteFriend(userId, groupId)
}

// 根据 群Id 查找所有群用户并 返回UsersList

func SearchUsersByGroupId(groupId string) []*UserBasic {
	var usersList []*UserBasic
	// 获取所有关系列表
	contactList = getContactList()

	for _, contact := range contactList {
		if contact.OwnerId == groupId && contact.Type == 2 {
			friendList = append(friendList, contact)
		}
	}
	// 将 []contactList 转化成 []UserBasic
	for _, item := range contactList {
		currentUser := FindUserByUID(item.TargetId)
		// 添加
		usersList = append(usersList, &currentUser)
	}
	return usersList
}
