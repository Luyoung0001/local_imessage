package models

import (
	"context"
	"encoding/json"
	"local_imessage/utils"
	"log"
)

type GroupBasic struct {
	Name       string   // 群名
	GroupID    string   // 群ID,由群主的 ID 生成
	OwnerUID   string   // 群主
	ManagerIDs []string // 管理员们
	Icon       string
	Desc       string
	Type       string // 保留字段,预留

}

func (table *GroupBasic) TableName() string {
	return "groupBasic"
}

// 创建群

func CreatGroup(userId string, group GroupBasic) {
	// group 字段有:
	// Name 等
	// Id 需要生成
	ctx := context.Background()
	// userId 创建了一个群聊,然后将整个群视为 UserBasic 处理
	// 将群序列化后存储
	// 获取 key
	groupId := utils.MD5Encode(userId)

	group.GroupID = groupId
	group.OwnerUID = userId
	// 序列化
	groupJSON, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}
	// 存储
	err = utils.Red.Set(ctx, groupId, groupJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

// 添加管理员

func AddMan(userId, groupId string) {
	ctx := context.Background()
	// 给 group 的字段增加新的对象,之后更新
	// 查询群
	utils.Red.Get(ctx, groupId).Result()

}
