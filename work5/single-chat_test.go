package main

import (
	"encoding/json"
	"five/config"
	"five/consts"
	"five/types"
	"fmt"
	"github.com/bytedance/gopkg/lang/fastrand"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"
)

var TestRabbitMQConn *amqp.Connection

const maxUserCount = consts.MaxUserCount
const cycleNumber = consts.CycleNumber

type loginResp struct {
	Data data        `json:"data"`
	Info interface{} `json:"info"`
}

type data struct {
	ID    int         `json:"id"`
	Token token       `json:"token"`
	Info  interface{} `json:"info"`
}

type token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func TestSingleChat(t *testing.T) {
	config.InitConfig()
	initRabbitMQ(t)
	// 注册并登录，获得用户信息。
	port := "8080"
	fmt.Println("开始执行")
	infos := userInfos(port, t)
	resp := make(chan []*websocket.Conn)
	go wsConns(infos, resp, port, t)
	conns := <-resp

	defer func() {
		for i := 0; i < maxUserCount; i++ {
			if conns[i] != nil {
				conns[i].Close()
			}
		}
	}()

	// 为了方便，将用户信息转换为map，方便通过uid查找对应的websocket连接。
	uids := make(map[int]int, maxUserCount) // uid -> index
	for i := 0; i < maxUserCount; i++ {
		uids[infos[i].Data.ID] = i
	}

	// 用于管理和通信consumer
	clientID := make(chan uint, maxUserCount)
	defer close(clientID)
	consumeDone := make(chan bool)
	defer close(consumeDone)
	go rabbitMQConsumer(clientID, consumeDone, t)

	msg := types.Message{
		Content: "hello",
		Type:    1,
	}

	f, err := os.Create("cpu.pprof")
	if err != nil {
		t.Fatal(err)
	}
	if err = pprof.StartCPUProfile(f); err != nil {
		t.Fatal(err)
	}

	start := time.Now().UnixMicro()
	defer pprof.StopCPUProfile()

	for i := 0; i < cycleNumber; i++ {
		uid := <-clientID
		conn := conns[uids[int(uid)]]
		conn.WriteJSON(msg)
	}
	end := time.Now().UnixMicro()
	fmt.Printf("发送%d条消息，耗时%d毫秒\n", cycleNumber, end-start)
	consumeDone <- true

	memFile, err := os.Create("mem.pprof")
	if err != nil {
		t.Fatal(err)
	}
	if err = pprof.WriteHeapProfile(memFile); err != nil {
		t.Fatal(err)
	}
	return
}

func userInfos(port string, t *testing.T) []loginResp {
	t.Helper()
	datas := make([]loginResp, maxUserCount)
	for i := 0; i < maxUserCount; i++ {
		randomNames := randomName()
		url := fmt.Sprintf("http://localhost:%s/api/user/register?user_name=%s&password=123456", port, randomNames[i])
		req, _ := http.NewRequest("POST", url, nil)
		client := &http.Client{}
		resp, _ := client.Do(req)
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}

		url = fmt.Sprintf("http://localhost:%s/api/user/login?user_name=%s&password=123456", port, randomNames[i])
		req, _ = http.NewRequest("POST", url, nil)
		resp, _ = client.Do(req)
		body, _ := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
		info := loginResp{}
		err = json.Unmarshal(body, &info)
		if err != nil {
			log.Println(err)
		}
		datas[i] = info
	}
	return datas
}

func wsConns(infos []loginResp, resp chan []*websocket.Conn, port string, t *testing.T) {
	t.Helper()
	conns := make([]*websocket.Conn, maxUserCount)
	for i := 0; i < maxUserCount; i++ {
		url := fmt.Sprintf("ws://localhost:%s/api/user/websocket?connect_type=1&receiver_id=%d", port, aReceiverID(i, infos))
		header := http.Header{}
		header.Add("access_token", infos[i].Data.Token.AccessToken)
		header.Add("refresh_token", infos[i].Data.Token.RefreshToken)
		c, _, err := websocket.DefaultDialer.Dial(url, header)
		if err != nil {
			t.Errorf("dialed error: %v\n", err)
		}
		conns[i] = c
	}
	resp <- conns
	return
}

func aReceiverID(i int, infos []loginResp) int {
	for {
		id := fastrand.Intn(maxUserCount)
		if id != i {
			return infos[id].Data.ID
		}
	}
}

func randomName() []string {
	names := make([]string, maxUserCount)
	for i := 0; i < maxUserCount; i++ {
		var result string
		for j := 0; j < 6; j++ {
			result += string(rune('a' + fastrand.Intn(26)))
		}
		names[i] = result
	}
	return names
}

func rabbitMQConsumer(clientID chan uint, done chan bool, t *testing.T) {
	t.Helper()
	conn := TestRabbitMQConn
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	msgs, err := ch.Consume(consts.SendDoneSingle, "", true, false, false, false, nil)
	if err != nil {
		t.Errorf("consume error: %v\n", err)
	}
	for {
		select {
		case msg := <-msgs:
			id, _ := strconv.ParseUint(string(msg.Body), 10, 32)
			clientID <- uint(id)
		case <-done:
			fmt.Println("consumer over")
			return
		}
	}
}

func initRabbitMQ(t *testing.T) {
	t.Helper()
	conf := config.Config.RabbitMQ
	url := conf.RabbitMQ + "://" + conf.RabbitMQUser + ":" + conf.RabbitMQPassWord + "@" + conf.RabbitMQHost + ":" + conf.RabbitMQPort + "/"
	conn, err := amqp.Dial(url)
	TestRabbitMQConn = conn
	if err != nil {
		t.Error(err)
	}
}
