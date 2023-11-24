package ws

import (
	"encoding/json"
	"five/consts"
	log2 "five/pkg/log"
	"five/pkg/myutils"
	"five/repository/db/dao"
	"five/repository/db/model"
	"five/repository/rabbitmq"
	"five/repository/redis"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

var allWsClient = make(map[uint]*Client)

func SingleChatStart(sManager ClientManager) {
	manager := &sManager
	for {
		select {
		case client := <-manager.Register:
			log.Printf("单聊建立新连接:%d->%d", client.UID, client.ReceiverID)
			allWsClient[client.UID] = client
			_ = client.Socket.WriteMessage(websocket.TextMessage, []byte("欢迎登录聊天室"))

		case client := <-manager.Unregister:
			log.Printf("单聊断开连接:%d->%d", client.UID, client.ReceiverID)
			if _, ok := allWsClient[client.UID]; ok {
				delete(allWsClient, client.UID)
				close(client.Message)
			}

		case broadcast := <-manager.Broadcast:
			message := broadcast.Message
			var readTag int8
			// 接收者的连接
			if conn, ok := allWsClient[broadcast.Client.ReceiverID]; ok {
				replyMsg := &SingleChatReplyMsg{
					From:    broadcast.Client.ReceiverID,
					Content: string(message),
				}
				msg, _ := json.Marshal(replyMsg)
				conn.Message <- msg
				readTag = consts.ReadMessage
			} else {
				log2.LogrusObj.Errorln("接收者不存在在该Manager中,发送者:", broadcast.Client.UID, "接收者:", broadcast.Client.ReceiverID)
				continue
			}

			msg := &model.Message{
				Uid:         broadcast.Client.UID,
				ReceiverID:  broadcast.Client.ReceiverID,
				Content:     string(message),
				MessageType: broadcast.MessageType,
				ReadTag:     readTag,
				CreateAt:    time.Now().Unix(),
			}

			rabbitmq.PublishMessage(msg)
		}
	}
}

func GroupChatStart(gManager ClientManager) {
	manager := &gManager
	clients := manager.Clients.(map[uint]map[uint]*Client)
	for {
		select {
		case conn := <-manager.Register: // 群聊建立新连接
			log.Printf("群聊建立新连接:%d->%d", conn.UID, conn.ReceiverID)
			if _, ok := clients[conn.ReceiverID]; !ok { // 如果没有这个群聊，就创建一个
				clients[conn.ReceiverID] = make(map[uint]*Client)
			}
			key := "group:" + strconv.FormatUint(uint64(conn.ReceiverID), 10)
			redis.RedisClient.SAdd(key, conn.UID) // 把用户加入到群聊中
			clients[conn.ReceiverID][conn.UID] = conn
			_ = conn.Socket.WriteMessage(websocket.TextMessage, []byte("欢迎登录聊天室"))

		case conn := <-manager.Unregister: // 群聊断开连接
			if _, ok := clients[conn.ReceiverID][conn.UID]; ok {
				delete(clients[conn.ReceiverID], conn.UID)
				close(conn.Message)
			}
			if len(clients[conn.ReceiverID]) == 0 {
				delete(clients[conn.ReceiverID], conn.ReceiverID)
			}

		case broadcast := <-manager.Broadcast: // 群聊消息广播
			message := broadcast.Message
			receiverID := broadcast.Client.ReceiverID
			onlineMember := clients[receiverID]
			// 获取群聊中的所有成员
			key := "group:" + myutils.UintToString(receiverID)
			members, _ := redis.RedisClient.SMembers(key).Result()
			// 遍历所有的成员，如果在线就发送消息，不在线就不发送
			for _, member := range members {
				mid := myutils.StringToUint(member)
				replyMsg := &GroupChatReplyMsg{ // 用于返回给发送者的消息
					From:    broadcast.Client.UID,
					GroupID: receiverID,
					Content: string(message),
				}
				msg := &model.Message{ // 用于保存到数据库的消息
					Uid:         broadcast.Client.UID,
					ReceiverID:  receiverID, // 群聊的接收者是群聊的ID
					Content:     string(message),
					MessageType: broadcast.MessageType,
					ReadTag:     consts.ReadMessage,
					CreateAt:    time.Now().Unix(),
				}
				rMsg, _ := json.Marshal(replyMsg) // 转换为json格式

				if mid != broadcast.Client.UID { // 不给自己发消息
					if conn, ok := onlineMember[mid]; ok { // 如果在线就发送，并标记为在线消息
						conn.Message <- rMsg
						msg.ReadTag = consts.ReadMessage
					} else { // 不在线就标记为离线消息
						msg.ReadTag = consts.UnReadMessage
					}
				}
				// 保存到数据库
				if err := dao.SaveGroupChatMessage(msg, mid); err != nil {
					log2.LogrusObj.Errorln(err)
					continue
				}
			}
		}
	}
}
