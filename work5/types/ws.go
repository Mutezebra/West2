package types

type WsConnectReq struct {
	ConnectType int8 `json:"connect_type,omitempty" form:"connect_type" binding:"required"` // 1的话表示单聊，2的话表示群聊
	ReceiverID  uint `json:"receiver_id,omitempty" form:"receiver_id" binding:"required"`   // 单聊的话表示用户ID，群聊的话就是群ID
}

type Message struct {
	Content   string `json:"content,omitempty"`
	Type      int8   `json:"type,omitempty"` // 1的话表示发送文本，2的话表示拉取历史消息，3的话表示搜索消息
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// GroupInfoResp 群组信息
type GroupInfoResp struct {
	ID        uint   `json:"id,omitempty"`
	GroupName string `json:"group_name,omitempty"`
}
