package e

var MsgFlags = map[int]string{
	SUCCESS: "Operator success",
	ERROR:   "Operator failed",

	ParseTokenFailed: "Parse token failed",
	CheckTokenFailed: "Check token failed",

	JsonUnmarshalFailed: "Json unmarshal failed",
	InvalidParams:       "Invalid params",

	// user
	UserExist:           "User exist",
	SetPasswordFailed:   "Set password failed",
	UserNotExist:        "User not exist",
	PasswordError:       "Password error",
	GenerateTokenFailed: "Generate token failed",
	GetUserInfoFailed:   "Get user info failed",

	// ws
	WebsocketSuccessMessage: "解析content内容信息",
	WebsocketSuccess:        "发送信息，请求历史纪录操作成功",
	WebsocketEnd:            "请求历史纪录，但没有更多记录了",
	WebsocketOnlineReply:    "针对回复信息在线应答成功",
	WebsocketOfflineReply:   "针对回复信息离线回答成功",
	WebsocketLimit:          "请求收到限制",
	WsConnectFailed:         "websocket connect failed",
	GroupNotExist:           "group not exist",
	WsConnectTypeFailed:     "websocket connect type failed",
	GroupExist:              "group exist",
	CreateGroupFailed:       "create group failed",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
