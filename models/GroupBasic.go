package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Desc    string
	Type    string // 保留字段,预留

}

func (table *GroupBasic) TableName() string {
	return "groupBasic"
}
