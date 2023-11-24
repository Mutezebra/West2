package service

import (
	"context"
	"five/consts"
	"five/pkg/ctl"
	"five/pkg/e"
	"five/pkg/log"
	"five/repository/db/dao"
	"five/service/ws"
	"five/types"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WsService struct{}

var WsSrv *WsService
var WsSrvOnce sync.Once
var upGrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func GetWsSrv() *WsService {
	WsSrvOnce.Do(func() {
		WsSrv = &WsService{}
	})
	return WsSrv
}

// 时间戳为index，根据时间戳向上查找10条信息。额外定制一个一定时间范围内的查询。
func (s *WsService) Connect(ctx context.Context, req *types.WsConnectReq, w http.ResponseWriter, r *http.Request) {
	var code = e.SUCCESS
	// 将原始HTTP连接升级为基于websocket的连接
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		code = e.WsConnectFailed
		log.LogrusObj.Errorln(err)
		return
	}

	// 根据连接类型，判断连接对象是否存在
	switch req.ConnectType {
	case consts.ConnectTypeSingle:
		userDao := dao.GetUserDao(ctx)
		_, err := userDao.FindUserByID(req.ReceiverID)
		if err != nil {
			code = e.UserNotExist
			log.LogrusObj.Errorln(e.GetMsg(code))
			conn.WriteJSON(ctl.Response{Status: code, Msg: e.GetMsg(code), Error: err.Error()})
			return
		}
	case consts.ConnectTypeGroup:
		groupDao := dao.GetGroupDao(ctx)
		_, err := groupDao.FindGroupByID(req.ReceiverID)
		if err != nil {
			code = e.GroupNotExist
			log.LogrusObj.Errorln(e.GetMsg(code))
			conn.WriteJSON(ctl.Response{Status: code, Msg: e.GetMsg(code), Error: err.Error()})
			return
		}
	default:
		code = e.WsConnectTypeFailed
		log.LogrusObj.Errorln(e.GetMsg(code))
		conn.WriteJSON(ctl.Response{Status: code, Msg: e.GetMsg(code), Error: "连接类型错误"})
		return
	}

	// 从上下文中获取用户信息
	userInfo := ctl.GetFromContext(ctx)
	if userInfo == nil {
		code = e.GetUserInfoFailed
		log.LogrusObj.Errorln(e.GetMsg(code))
		conn.WriteJSON(ctl.Response{Status: code, Msg: e.GetMsg(code), Error: "用户信息获取失败"})
		return
	}

	// 实例化一个客户端
	client := &ws.Client{
		UID:        userInfo.UID,
		ReceiverID: req.ReceiverID,
		ClientType: req.ConnectType,
		Socket:     conn,
		Message:    make(chan []byte),
	}

	// 启动协程，监听用户发送的消息
	go client.Read()
	go client.Write()
}
