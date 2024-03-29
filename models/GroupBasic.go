package models

import (
	"context"
	"encoding/json"
	"fmt"
	"local_imessage/utils"
	"log"
)

type GroupBasic struct {
	Name        string   // 群名
	GroupID     string   // 群ID,由群主的 ID 生成
	OwnerUID    string   // 群主
	ManagerIDs  []string // 管理员们
	Description string   // 描述
	DataType    string   `json:"dataType"` // 添加 DataType 字段来标识数据类型
}

func (table *GroupBasic) TableName() string {
	return "groupBasic"
}

// 查看所有的群

func GetGroupList() []GroupBasic {
	groupList := make([]GroupBasic, 0)
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
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		if StructType(value) == "groupBasic" {
			var groupInfo GroupBasic
			// 反序列化
			err = json.Unmarshal([]byte(value), &groupInfo)
			if err != nil {
				return nil
			}
			// 添加到 values
			groupList = append(groupList, groupInfo)
		}
	}
	return groupList
}

// 创建群
// 需要注意的是,创建群也会产生一个 contact 类型的记录

func CreatGroup(group GroupBasic) bool {
	ctx := context.Background()
	group.GroupID = IdGenerator()
	group.DataType = "groupBasic"

	// 序列化
	groupJSON, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 存储
	err = utils.Red.Set(ctx, group.GroupID, groupJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 加入群
	return JoinGroup(group.OwnerUID, group.GroupID)

}

// 删除群聊

func DeleteGroup(operator, groupId string) bool {
	ctx := context.Background()
	// 鉴定身份
	currentGroup := FindGroupByGID(groupId)

	if currentGroup.OwnerUID != operator {
		return false
	} else {
		// 要删除的用户键
		userKey := groupId

		deleted, err := utils.Red.Del(ctx, userKey).Result()
		if err != nil {
			log.Println("Error:", err)
			return false
		}
		if deleted > 0 {
			fmt.Printf("Group with key '%s' deleted.\n", userKey)
		} else {
			fmt.Printf("Group with key '%s' not found.\n", userKey)
			return false
		}
	}
	return true
}

// 添加管理员

func AddMan(userId, groupId string) bool {
	ctx := context.Background()
	// 给 group 的字段增加新的对象,之后更新
	// 查询群
	groupRaw, err := utils.Red.Get(ctx, groupId).Result()
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 反序列化
	var currentGroup GroupBasic
	err = json.Unmarshal([]byte(groupRaw), &currentGroup)
	// 修改
	currentGroup.ManagerIDs = append(currentGroup.ManagerIDs, userId)
	// 序列化
	currentJSON, err := json.Marshal(currentGroup)
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 重新存储
	err = utils.Red.Set(ctx, currentGroup.GroupID, currentJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// 通过群 ID 找 group

func FindGroupByGID(gid string) GroupBasic {
	ctx := context.Background()
	var group GroupBasic
	groupJSON, err := utils.Red.Get(ctx, gid).Result()
	if err != nil {
		log.Println(err)
		return GroupBasic{}
	}
	err = json.Unmarshal([]byte(groupJSON), &group)
	if err != nil {
		return GroupBasic{}
	}
	return group
}

// 鉴定权限

func LevelUserInGroup(operatorId, groupId, userId string) bool {
	ctx := context.Background()
	group := FindGroupByGID(groupId)
	// 谁能删除谁?
	// 1. 群主可以删除任何人,想要删除自己,群主可以直接删除群聊
	// 2. 管理员仅仅可以删除普通成员,想要删除自己,可以自己退群
	// 3. 普通成员不可以删除任何人,想要删除自己,可以自己退群

	// 鉴定权限
	// 1.operator 是那一个级别?
	// 1:群主,0:管理员,-1:普通成员
	var leverOp int
	if group.OwnerUID == operatorId {
		leverOp = 1
	} else if stringInSlice(operatorId, group.ManagerIDs) {
		leverOp = 0
	} else {
		leverOp = -1
	}

	// 2.被删除对象 userId 是哪一个级别?
	var leverUser int
	if group.OwnerUID == userId {
		leverUser = 1
	} else if stringInSlice(userId, group.ManagerIDs) {
		leverUser = 0
	} else {
		leverUser = -1
	}
	// 比较
	if leverOp > leverUser {
		// 删除群的成员
		// 调用 DeleteFriend()
		DeleteFriend(userId, groupId)
		// 如果是管理员,继续修改管理员列表
		if leverUser == 1 {
			// 修改 group.ManagerIDs 字段
			group.ManagerIDs = removeFromSliceUsingCopy(group.ManagerIDs, userId)
			// 序列化
			groupJSON, err := json.Marshal(group)
			if err != nil {
				log.Fatal(err)
				return false
			}
			// 存储
			utils.Red.Set(ctx, group.GroupID, groupJSON, 0)
		}
		return true
	} else {
		return false
	}

}

// 判断 userId 和 GroupId 之间的关系
// 2:群主;1:管理员;0:普通成员;-1:不是群员

func RelationBetweenUserAndGroup(userId, groupId string) int {
	var relation int
	currentGroup := FindGroupByGID(groupId)
	if currentGroup.GroupID == userId {
		relation = 2
	} else if stringInSlice(userId, currentGroup.ManagerIDs) {
		relation = 1
	} else if ContactRelation(userId, groupId) > 0 {
		relation = 0
	} else {
		relation = -1
	}
	return relation
}
