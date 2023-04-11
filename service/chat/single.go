package chat

import (
	"chat/cache"
	"chat/model"
	"chat/pkg/e"
	"chat/pkg/utils"
	"chat/serializer"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

// Start 单聊初始化
func (manager *Manager) Start() {
	for {
		log.Println("<---监听管道通信--->")
		select {

		case conn := <-ClientsManager.Login: // 建立连接
			log.Printf("建立新连接: %v", conn.ID)
			ClientsManager.Clients[conn.ID] = conn
			replyMsg := &SingleMessage{
				FromID:     "0",
				FromUser:   "服务器",
				Code:       e.WebsocketSuccess,
				Content:    "已连接至服务器",
				CreateTime: utils.GetLocalDateTime(),
			}
			msg, _ := json.Marshal(replyMsg)
			err := conn.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("InsertOneMsg Err", err)
			}

		case conn := <-ClientsManager.UnLogin: // 断开连接
			log.Printf("连接失败:%v", conn.ID)

			if _, ok := ClientsManager.Clients[conn.ID]; ok {

				replyMsg := &SingleMessage{
					FromID:     "0",
					FromUser:   "服务器",
					Code:       e.WebsocketEnd,
					Content:    "连接已断开",
					CreateTime: utils.GetLocalDateTime(),
				}

				msg, _ := json.Marshal(replyMsg)
				err := conn.Socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
				close(conn.Send)
				delete(ClientsManager.Clients, conn.ID)
			}

		case broadcast := <-ClientsManager.Broadcast: //广播消息
			message := broadcast.Message
			sendId := broadcast.User.SendID
			flag := false // 默认对方不在线

			for id, conn := range ClientsManager.Clients {
				if id != sendId {
					continue
				}
				select {
				// 发送消息
				case conn.Send <- message:
					flag = true
				default:
					// 如果发送失败，就关闭这个连接
					close(conn.Send)
					delete(ClientsManager.Clients, conn.ID)
				}
			}
			if flag {
				log.Println("对方在线应答")

				replyMsg := &SingleMessage{
					FromID:     "0",
					FromUser:   "服务器",
					Code:       e.WebsocketOnlineReply,
					Content:    "对方在线应答",
					CreateTime: utils.GetLocalDateTime(),
				}

				msg, err := json.Marshal(replyMsg)
				_ = broadcast.User.Socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			} else {
				log.Println("对方不在线")

				replyMsg := SingleMessage{
					FromID:     "0",
					FromUser:   "服务器",
					Code:       e.WebsocketOfflineReply,
					Content:    "对方不在线应答",
					CreateTime: utils.GetLocalDateTime(),
				}

				msg, err := json.Marshal(replyMsg)
				err = broadcast.User.Socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			}
		}
	}
}

// SingleChat 单聊
func SingleChat(c *gin.Context) {
	//chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	//fmt.Println("Authorization", chaim.Id)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议

	if err != nil {
		fmt.Println(err)
		http.NotFound(c.Writer, c.Request)
		return
	}
	chaim, _ := utils.ParseToken(c.GetHeader("Sec-WebSocket-Protocol"))

	// 创建一个用户实例
	client := &Client{
		ID:       strconv.Itoa(int(chaim.Id)),
		UserName: chaim.UserName,
		Socket:   conn,
		Send:     make(chan SingleMessage),
	}

	// 用户注册到用户管理上
	ClientsManager.Login <- client
	go client.Read()
	go client.Write()
}

// 读取消息
func (c *Client) Read() {
	defer func() { // 避免忘记关闭，所以要加上close
		ClientsManager.UnLogin <- c
		_ = c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		sendMsg := new(SingleSendMsg)
		// _,msg,_:=c.Socket.ReadMessage()
		err := c.Socket.ReadJSON(&sendMsg) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			log.Println("数据格式不正确", err)
			ClientsManager.UnLogin <- c
			_ = c.Socket.Close()
			break
		}
		//发送消息
		if sendMsg.Type == 1 {

			if len(sendMsg.Content) == 0 {
				continue
			}

			log.Println(c.ID, "发送消息", sendMsg.Content)

			if sendMsg.ToUid == "0" {
				continue
			}

			c.SendID = sendMsg.ToUid

			message := SingleMessage{
				FromID:     c.ID,
				FromUser:   c.UserName,
				ToID:       c.SendID,
				Code:       e.WebsocketSuccessMessage,
				Content:    sendMsg.Content,
				CreateTime: utils.GetLocalDateTime(),
				ExpireTime: utils.GetExpireDateTime(),
			}

			ClientsManager.Broadcast <- &Broadcast{
				User:    c,
				Message: message,
			}
			//go singleMysqlSave(message)
			//储存
			go singleRedisSave(message)
			//获取聊天记录
		} else if sendMsg.Type == 2 {

			if sendMsg.ToUid == "0" {
				continue
			}

			fromId, _ := strconv.ParseInt(c.ID, 10, 64)
			toId, _ := strconv.ParseInt(sendMsg.ToUid, 10, 64)
			name := ""
			if toId < fromId {
				name = fmt.Sprintf("%s-%s", sendMsg.ToUid, c.ID)
			} else {
				name = fmt.Sprintf("%s-%s", c.ID, sendMsg.ToUid)
			}

			messages := cache.RedisClient.LRange(context.TODO(), name, 0, -1).Val()

			var singleMessages []model.SingleMessage
			for _, message := range messages {

				var singleMessage model.SingleMessage
				json.Unmarshal([]byte(message), &singleMessage)

				//过期消息删除
				t1, err1 := utils.TranslateTime(singleMessage.ExpireTime)
				if err1 != nil {
					fmt.Println(err1)
				}
				t2 := utils.GetLocalTime()
				if t2.After(t1) {
					cache.RedisClient.LRem(context.TODO(), name, 1, message)
					continue
				}

				singleMessages = append(singleMessages, singleMessage)
			}

			content, _ := json.Marshal(serializer.BuildSingleMessages(singleMessages))
			replyMsg := SingleMessage{
				FromID:     "0",
				FromUser:   "服务器",
				Code:       e.WebsocketSuccessMessage,
				Content:    string(content),
				CreateTime: utils.GetLocalDateTime(),
			}

			msg, err1 := json.Marshal(replyMsg)
			if err1 != nil {
				fmt.Println("json格式转换失败", err)
			}

			err = c.Socket.WriteMessage(websocket.TextMessage, msg)
			if err1 != nil {
				fmt.Println("InsertOneMsg Err", err)
			}
		}
	}
}

// 发送消息
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {

		select {

		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println(c.ID, "接受消息:", message.Content)

			msg, err := json.Marshal(message)
			if err != nil {
				fmt.Println("json格式转换失败", err)
			}

			err = c.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("InsertOneMsg Err", err)
			}
		}
	}
}
func singleMysqlSave(replyMsg SingleMessage) {
	var msg model.SingleMessage

	msg.MessageData = replyMsg.Content
	msg.CreateTime = replyMsg.CreateTime
	msg.ToId, _ = strconv.ParseInt(replyMsg.ToID, 10, 64)
	msg.FromId, _ = strconv.ParseInt(replyMsg.FromID, 10, 64)

	err := model.DB.Save(&msg).Error
	if err != nil {
		fmt.Println("消息mysql存储失败！！！")
	}
}

// redis储存
func singleRedisSave(replyMsg SingleMessage) {
	var msg model.SingleMessage

	msg.MessageData = replyMsg.Content
	msg.CreateTime = replyMsg.CreateTime
	msg.ToId, _ = strconv.ParseInt(replyMsg.ToID, 10, 64)
	msg.FromId, _ = strconv.ParseInt(replyMsg.FromID, 10, 64)
	msg.ExpireTime = replyMsg.ExpireTime

	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json格式转换失败", err)
	}

	name := ""
	if replyMsg.ToID < replyMsg.FromID {
		name = fmt.Sprintf("%s-%s", replyMsg.ToID, replyMsg.FromID)
	} else {
		name = fmt.Sprintf("%s-%s", replyMsg.FromID, replyMsg.ToID)
	}

	_, err = cache.RedisClient.RPush(context.TODO(), name, string(message)).Result()
	if err != nil {
		fmt.Println("消息mysql存储失败！！！")
	}
	//发送消息
	cache.SingleStreamMQ.SendMsg(context.Background(), &cache.Msg{
		Topic: "single",
		Body:  message,
	})
}
