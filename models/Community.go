package models

import (
	"fmt"
	"gorm.io/gorm"
	"local_imessage/utils"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

// [禁止] 使用 panic 用于正常的错误处理。
// 应该使用 error 和多个返回值。
// panic 只适合用于那些严重影响程序运行的错误

func CreateCommunity(community Community) (int, string) {
	tx := utils.DB.Begin()
	//defer func() {
	//	// recover():内建函数;函数用于从发生的恐慌（panic）中恢复
	//	if r := recover(); r != nil {
	//		// 事物回滚
	//		tx.Rollback()
	//	}
	//}()

	if len(community.Name) == 0 {
		return -1, "群名称不能为空!"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录!"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return -1, "建群失败!"
	}
	contact := Contact{}
	contact.OwnerId = community.OwnerId
	contact.TargetId = community.ID
	contact.Type = 2 //群关系
	if err := utils.DB.Create(&contact).Error; err != nil {
		tx.Rollback()
		return -1, "添加群关系失败!"
	}
	// 提交事务
	tx.Commit()
	return 0, "建群成功!"

}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=2", ownerId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}

	data := make([]*Community, 10)
	utils.DB.Where("id in ?", objIds).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data, "查询成功"
}
