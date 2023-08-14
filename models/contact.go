package models

import (
	"context"
	"encoding/json"
	"local_imessage/utils"
	"log"
)

type Contact struct {
	ContactID string // 哈哈哈哈哈哈哈哈笑会儿
	OwnerId   string // userId
	TargetId  string // 有人的 ID,有群的 ID
	Type      int    // 1 代表好友关系;2 代表群聊关系
	Desc      string // 描述? 预留字段
	DataType  string `json:"dataType"` // 添加 DataType 字段来标识数据类型
}

func (table *Contact) TableName() string {
	return "contact"
}

// 返回所有的列表,这个列表为后面的查询等提供服务

func GetContactList() []Contact {
	contactList := make([]Contact, 0)
	ctx := context.Background()
	var keys []string
	// 获取所有键
	keys, err := utils.Red.Keys(ctx, "*").Result()
	if err != nil {
		log.Fatal(err)
	}
	// 遍历每个键，获取值,并对每个值进行判断,看是否是 contact
	for _, key := range keys {
		value, err := utils.Red.Get(ctx, key).Result()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		// 判断数据类型
		if StructType(value) == "contact" {
			var contact Contact
			err = json.Unmarshal([]byte(value), &contact)
			if err != nil {
				log.Println("Error:", err)
				continue
			}
			contactList = append(contactList, contact)
		}
	}
	return contactList

}

// 返回 userId 的所有好友列表

func FriendList(userId string) []Contact {
	contactList := GetContactList()
	friendList := make([]Contact, 0) // 创建一个空的好友列表

	for _, contact := range contactList {
		if contact.OwnerId == userId && contact.Type == 1 {
			// 创建一个新的 Contact 对象，并将属性拷贝过来
			newContact := Contact{
				ContactID: contact.ContactID,
				OwnerId:   contact.OwnerId,
				TargetId:  contact.TargetId,
				Type:      contact.Type,
				Desc:      "",
				DataType:  "contact",
			}
			friendList = append(friendList, newContact)
		}
	}

	return friendList
}

// 添加好友

func AddFriend(userId string, targetID string) bool {
	return CreateRelation(userId, targetID, 1)
}

// userId 加入 groupId 的群聊

func JoinGroup(userId string, groupId string) bool {

	return CreateRelation(userId, groupId, 2)
}

// 删除好友
// 删除两个contact 记录

func DeleteFriend(userId, targetId string) bool {
	//删除的前提是对方已经成为你的好友,这个由上层控制,这里不做校验
	contact1 := FindContactByOwnerIdAndTargetId(userId, targetId)
	contact2 := FindContactByOwnerIdAndTargetId(targetId, userId)
	return DeleteKvFromRed(contact1.ContactID) && DeleteKvFromRed(contact2.ContactID)
}

// 找到 ownerId 为 userId,且同时 targetId 为targetId的 contact 记录

func FindContactByOwnerIdAndTargetId(userId, targetId string) Contact {
	contacts := FriendList(userId)
	// 在 contacts 中找targetId 为targetId的 contact 记录
	for _, each := range contacts {
		if each.TargetId == targetId {
			return each
		}
	}
	return Contact{}

}

// 退群
// 删除两个 contact 记录

func OutGroup(userId, groupId string) bool {
	// 双相删除记录
	return DeleteFriend(userId, groupId)
}

// 返回用户 userID 加的所有群,包括自己创建的群

func GroupList(userID string) []Contact {
	var userGroupList []Contact
	// 就是查找type=2
	contactList := FriendList(userID)
	for _, item := range contactList {
		if item.Type == 2 {
			userGroupList = append(userGroupList, item)
		}
	}
	return userGroupList
}

// 根据 群Id 查找所有群用户并 返回UsersList
// 相当于找 groupId 的好友

func SearchUsersByGroupId(groupId string) []UserBasic {

	usersList := make([]UserBasic, 0)
	targets := make([]string, 0)
	groupContact := FriendList(groupId)
	// 拿到所有的 []targetId
	for _, each := range groupContact {
		targets = append(targets, each.TargetId)
	}
	// 遍历所有的 targetId,找到 UserBasic
	for _, each := range targets {
		userBasic := FindUserByUID(each)
		usersList = append(usersList, userBasic)
	}
	return usersList
}

// 返回用户 userID 加的所有群 []GroupBasic

func GroupsList(userID string) []GroupBasic {
	var userGroupList []Contact
	var groupsList []GroupBasic
	// 加的所有的群 contacts
	userGroupList = GroupList(userID)
	// 匹配
	for _, each := range userGroupList {
		group := FindGroupByGID(each.TargetId)
		groupsList = append(groupsList, group)
	}
	return groupsList
}

// 判断两个成员是否有关系
// 1:好友关系;2:群和成员关系;-1:没有关系

func ContactRelation(userId1, userId2 string) int {

	// 获取该用户的所有好友
	var relation int
	contactList := GetContactList()
	for _, each := range contactList {
		if each.TargetId == userId2 && each.Type == 1 {
			relation = 1
		} else if each.TargetId == userId2 && each.Type == 2 {
			relation = 2
		} else {
			// 没有关系
			relation = -1
		}
	}
	return relation
}

// 判断contactId 是群记录还是用户记录
// 是群就返回 2
// 是用户就返回 1
// 错误返回 0

func ContactType(contactId string) int {
	ctx := context.Background()
	contactJSON, err := utils.Red.Get(ctx, contactId).Result()
	if err != nil {
		log.Println(err)
		return 0
	}
	// 解析
	var contact Contact
	err = json.Unmarshal([]byte(contactJSON), &contact)
	if err != nil {
		return 0
	}
	return contact.Type
}

func CreateRelation(userId string, targetID string, tp int) bool {
	ctx := context.Background()
	// 添加好友的本质就是创建一个contact 记录
	var contact1 Contact
	var contact2 Contact
	// 记录 1
	contact1.ContactID = IdGenerator()
	contact1.Type = tp // 关系
	contact1.OwnerId = userId
	contact1.TargetId = targetID
	contact1.DataType = "contact"
	//记录 2
	contact2.ContactID = IdGenerator()
	contact2.Type = tp // 关系
	contact2.OwnerId = targetID
	contact2.TargetId = userId
	contact2.DataType = "contact"

	// 序列化
	userJSON1, err1 := json.Marshal(contact1)
	userJSON2, err1 := json.Marshal(contact2)
	if err1 != nil {
		log.Fatal(err1)
		return false
	}
	// 写入 redis 并返回
	err1 = utils.Red.Set(ctx, contact1.ContactID, userJSON1, 0).Err()
	err1 = utils.Red.Set(ctx, contact2.ContactID, userJSON2, 0).Err()
	if err1 != nil {
		log.Fatal(err1)
		return false
	}
	return true

}

// 删除 contact 记录

func DeleteKvFromRed(contactId string) bool {
	ctx := context.Background()
	// 删除用户
	_, err := utils.Red.Del(ctx, contactId).Result()
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	return true
}

// 已知
