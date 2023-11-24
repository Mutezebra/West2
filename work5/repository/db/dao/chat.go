package dao

import (
	"context"
	"five/consts"
	"five/repository/db/model"
	"gorm.io/gorm"
)

type ChatDao struct {
	*gorm.DB
}

func GetChatDao(ctx context.Context) *ChatDao {
	return &ChatDao{NewDBClient(ctx)}
}

// SaveSingleChatMessage 保存单聊消息
func SaveSingleChatMessage(msg *model.Message) (err error) {
	db := NewDBClient(context.TODO())
	err = db.Model(&model.Message{}).Create(&msg).Error
	return
}

// GetChatHistoryMsg 获取聊天记录,一次获取十条
func GetChatHistoryMsg(uid, receiverID uint, createAt int64) (msgs []*model.Message, err error, existUnReadMsg bool) {
	db := NewDBClient(context.Background())
	if err = db.Model(&model.Message{}).
		Where("uid=? and receiver_id=? and create_at<?", uid, receiverID, createAt).
		Order("id desc").
		Limit(11).Find(&msgs).
		Error; err != nil {
		return nil, err, false
	}
	if len(msgs) > 10 {
		msgs = msgs[:10]
		existUnReadMsg = true
	}
	for _, msg := range msgs {
		if msg.ReadTag == consts.UnReadMessage {
			db.Model(&model.Message{}).Where("id=?", msg.ID).UpdateColumn("read_tag", consts.ReadMessage)
		}
	}
	return
}

// SearchChatHistoryMsg 搜索聊天记录
func SearchChatHistoryMsg(uid, receiverID uint, startTime, endTime int64) (msgs []*model.Message, err error) {
	db := NewDBClient(context.Background())
	if err = db.Model(&model.Message{}).
		Where("uid=? and receiver_id=? and create_at between ? and ?", uid, receiverID, startTime, endTime).
		Order("id desc").
		Find(&msgs).
		Error; err != nil {
		return nil, err
	}
	return
}
