package model

// Message 一条聊天消息的基本属性

type Message struct {
	ID          uint `gorm:"primarykey"`
	Uid         uint
	ReceiverID  uint
	Content     string
	MessageType int8  `gorm:"unsigned"` // 1表示文本消息
	ReadTag     int8  `gorm:"unsigned"` // 0表示未读，1表示已读
	CreateAt    int64 `gorm:"autoCreateTime,index"`
}
