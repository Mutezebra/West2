package model

import "gorm.io/gorm"

// Chat 一条信息的model，包含有所有应该具有的属性
type Chat struct {
	gorm.Model
	UserId    uint   `json:"user_id" gorm:"column:user_id;not null"`
	Content   string `json:"content" gorm:"column:content;not null"`
	ChatType  int8   `json:"chat_type" gorm:"column:chat_type;not null"`
	ChatState int8   `json:"chat_state" gorm:"column:chat_state;not null"`
}
