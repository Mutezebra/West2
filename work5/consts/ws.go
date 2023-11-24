package consts

// ConnectType 建立连接时的连接类型
const (
	ConnectTypeSingle = 1
	ConnectTypeGroup  = 2
)

// 用户发送的消息类型
const (
	MsgTypeText    = 1
	MsgTypeHistory = 2
	MsgTypeSearch  = 3
)

// ReadTag 消息的读取状态
const (
	UnReadMessage = 0
	ReadMessage   = 1
)

const (
	// ASyncWriteMysqlNumber 异步写入mysql的协程数
	ASyncWriteMysqlNumber = 10
	// SingleChatManagerNumber 单聊管理器的数量
	SingleChatManagerNumber = 5
	// GroupChatManagerNumber 群聊管理器的数量
	GroupChatManagerNumber = 5
)
