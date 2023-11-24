package ws

import (
	"encoding/json"
	"five/consts"
	"five/pkg/log"
	"five/pkg/myutils"
	"five/repository/db/dao"
	"five/repository/db/model"
	"five/repository/rabbitmq"
	"five/types"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"time"
)

type Client struct {
	UID        uint
	ReceiverID uint
	ClientType int8 // 1的话表示单聊，2的话表示群聊
	Message    chan []byte
	Socket     *websocket.Conn
}

type SingleChatReplyMsg struct {
	From    uint   `json:"from"`
	Content string `json:"content"`
}

type GroupChatReplyMsg struct {
	From    uint   `json:"from"`
	Content string `json:"content"`
	GroupID uint   `json:"group"`
}

type Broadcast struct {
	Client      *Client
	Message     []byte
	MessageType int8
}

type ClientManager struct {
	Clients    interface{}
	Broadcast  chan *Broadcast
	Register   chan *Client
	Unregister chan *Client
}

func (c *Client) Read() {
	manager := GetChatManager(c.ClientType)
	manager.Register <- c
	defer func() {
		fmt.Println("关闭Read连接")
		manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		// 读取客户端发送的消息
		msg := &types.Message{}
		_, msgBuf, err := c.Socket.ReadMessage()
		// 如果传过来的消息大于2048字节，就不处理
		if len(msgBuf) > 2048 {
			c.Socket.WriteMessage(websocket.TextMessage, []byte("消息过长"))
			continue
		}
		if err != nil && err == io.EOF {
			log.LogrusObj.Errorln("连接已关闭", err)
			break
		} else if err != nil {
			log.LogrusObj.Errorf("Socket ReadMessage failed，Err:%v", err)
			break
		}
		err = json.Unmarshal(msgBuf, &msg)
		if err != nil {
			msg.Type = 127
		}
		switch msg.Type { // 根据消息类型，做不同的处理
		case consts.MsgTypeText: // 文本消息
			manager.Broadcast <- &Broadcast{
				Client:      c,
				Message:     []byte(msg.Content),
				MessageType: consts.MsgTypeText,
			}
		case consts.MsgTypeHistory: // 历史消息
			var ExistUnReadMsg bool
			msgs := make([]*model.Message, 0)
			if c.ClientType == consts.ConnectTypeSingle {
				if msgs, err, ExistUnReadMsg = dao.GetChatHistoryMsg(c.ReceiverID, c.UID, time.Now().Unix()); err != nil {
					log.LogrusObj.Errorln(err)
					continue
				}
			} else if c.ClientType == consts.ConnectTypeGroup {
				// 此处的receiverID是groupID
				if msgs, err, ExistUnReadMsg = dao.GetGroupChatHistoryMsg(c.UID, c.ReceiverID, time.Now().Unix()); err != nil {
					log.LogrusObj.Errorln(err)
					continue
				}
			}
			if ExistUnReadMsg {
				c.Message <- []byte("仍有未读消息")
			} else {
				c.Message <- []byte("已无未读消息")
			}
			for _, item := range msgs {
				buf := make([]byte, 0)
				if buf, err = json.Marshal(item); err != nil {
					log.LogrusObj.Errorln(err)
					continue
				}
				c.Message <- buf
			}
		case consts.MsgTypeSearch:
			start, end, err := myutils.ParseTime(msg.StartTime, msg.EndTime)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.Socket.WriteMessage(websocket.TextMessage, []byte("时间格式不正确"))
				continue
			}
			var msgs []*model.Message
			if c.ClientType == consts.ConnectTypeSingle {
				if msgs, err = dao.SearchChatHistoryMsg(c.ReceiverID, c.UID, start, end); err != nil {
					log.LogrusObj.Errorln(err)
					continue
				}
			} else if c.ClientType == consts.ConnectTypeGroup {
				if msgs, err = dao.SearchGroupChatHistoryMsg(c.ReceiverID, start, end); err != nil {
					log.LogrusObj.Errorln(err)
					continue
				}
			}
			for _, item := range msgs {
				buf := make([]byte, 0)
				if buf, err = json.Marshal(item); err != nil {
					log.LogrusObj.Errorln(err)
					continue
				}
				c.Message <- buf
			}
		default:
			fmt.Println("msg不符合格式", msg)
			continue
		}
	}
}

func (c *Client) Write() {
	var err error
	uidStr := fmt.Sprintf("%d", c.UID)
	if ISTest {
		rabbitmq.SendSingle <- uidStr
	}
	defer func() {
		fmt.Println("关闭Write连接")
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Message: // 从管道中读取数据,管道中传过来的数据都是广播出来的
			if !ok { // 如果管道中没有数据，就关闭连接
				fmt.Println("管道中没有数据，关闭连接")
				if err = c.Socket.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.LogrusObj.Errorln(err)
				}
				return
			} // 否则就发送数据
			if err = c.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
				log.LogrusObj.Errorln(err)
				continue
			}
			if ISTest {
				rabbitmq.SendSingle <- uidStr
			}
		}
	}
}
