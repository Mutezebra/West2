package dao

import (
	"context"
	"five/consts"
	"five/pkg/myutils"
	"five/repository/db/model"
	"gorm.io/gorm"
)

type GroupDao struct {
	*gorm.DB
}

func GetGroupDao(ctx context.Context) *GroupDao {
	return &GroupDao{NewDBClient(ctx)}
}

func (dao *GroupDao) FindGroupByID(id uint) (group *model.Group, err error) {
	err = dao.DB.Model(&model.Group{}).First(&group, "id=?", id).Error
	return
}

// CreateTableByGroupID 根据GroupID创建一张表,表名为group_msg_groupID,用于存储群聊消息,表结构与message表相似，但是不需要其中的receiver_id字段
func (dao *GroupDao) CreateTableByGroupID(groupID uint) (err error) {
	err = dao.DB.Exec(
		"CREATE TABLE IF NOT EXISTS `group_msg_" + myutils.UintToString(groupID) +
			"` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`uid` bigint(20) unsigned NOT NULL," +
			"`member_id` bigint(20) unsigned NOT NULL," +
			"`content` varchar(255) NOT NULL," +
			"`message_type` tinyint(4) NOT NULL," +
			"`read_tag` tinyint(4) NOT NULL," +
			"`create_at` bigint(20) unsigned NOT NULL," +
			"PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;").
		Error
	return
}

// FindGroupByName 根据群聊名称查找群聊
func (dao *GroupDao) FindGroupByName(name string) (group *model.Group, err error) {
	err = dao.DB.Model(&model.Group{}).First(&group, "group_name=?", name).Error
	return
}

// CreateGroup 创建群聊
func (dao *GroupDao) CreateGroup(group *model.Group) (err error) {
	err = dao.DB.Model(&model.Group{}).Create(&group).Error
	if err != nil {
		return err
	}
	err = dao.CreateTableByGroupID(group.ID)
	return
}

// SaveGroupChatMessage 以msg中的group_msg_receiverID为表名，保存消息
func SaveGroupChatMessage(msg *model.Message, memberID uint) (err error) {
	db := NewDBClient(context.TODO())
	err = db.Exec(
		"INSERT INTO `group_msg_"+myutils.UintToString(msg.ReceiverID)+
			"` (`uid`, `member_id`, `content`, `message_type`, `read_tag`, `create_at`) VALUES (?, ?, ?, ?, ?, ?)",
		msg.Uid, memberID, msg.Content, msg.MessageType, msg.ReadTag, msg.CreateAt).
		Error
	return err
}

// GetGroupChatHistoryMsg 在表名为group_msg_groupID的表中，获取createID之前的十条消息
func GetGroupChatHistoryMsg(memberID, groupID uint, createID int64) (msgs []*model.Message, err error, existUnReadMsg bool) {
	db := NewDBClient(context.Background())
	if err = db.Table("group_msg_"+myutils.UintToString(groupID)).
		Where("create_at<? and member_id=?", createID, memberID).
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
			msg.ReadTag = consts.ReadMessage
		}
	}
	return
}

// SearchGroupChatHistoryMsg 搜索群聊消息
func SearchGroupChatHistoryMsg(groupID uint, startTime, endTime int64) (msgs []*model.Message, err error) {
	db := NewDBClient(context.Background())
	if err = db.Table("group_msg_"+myutils.UintToString(groupID)).
		Where("create_at between ? and ?", startTime, endTime).
		Order("id desc").
		Find(&msgs).
		Error; err != nil {
		return nil, err
	}
	return
}
