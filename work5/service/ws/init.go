package ws

import (
	"five/consts"
	"github.com/bytedance/gopkg/lang/fastrand"
)

var SingleChatManagers = make([]*ClientManager, consts.SingleChatManagerNumber, consts.SingleChatManagerNumber)
var GroupChatManager = ClientManager{
	Clients:    make(map[uint]map[uint]*Client),
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

// ISTest 为True的话说明在进行测试
var ISTest = false

func InitWsService() {
	for i := 0; i < consts.SingleChatManagerNumber; i++ {
		SingleChatManagers[i] = &ClientManager{
			Clients:    make(map[uint]*Client),
			Broadcast:  make(chan *Broadcast),
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
		}
	}
	for i := 0; i < consts.SingleChatManagerNumber; i++ {
		func(i int) {
			go SingleChatStart(*SingleChatManagers[i])
		}(i)
	}

	go GroupChatStart(GroupChatManager)
}

func TestInit() {
	ISTest = true
}

func GetChatManager(clientType int8) *ClientManager {
	if clientType == 1 {
		return SingleChatManagers[fastrand.Intn(consts.SingleChatManagerNumber)]
	} else {
		return &GroupChatManager
	}
}
