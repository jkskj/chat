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

// GroupChat 群聊
func GroupChat(c *gin.Context) {
	//chaim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	//升级get请求为webSocket协议
	ws, err := (&websocket.Upgrader{
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
	// 获取连接的用户id与群组id
	chaim, _ := utils.ParseToken(c.GetHeader("Sec-WebSocket-Protocol"))
	userId := strconv.Itoa(int(chaim.Id))
	groupId := c.Query("gid")

	//将id转为int64类型
	gid, _ := strconv.ParseInt(groupId, 10, 64)
	uid, _ := strconv.ParseInt(userId, 10, 64)

	//将当前连接的客户端放入客户端map（clients）中
	wsKey := WsKey{
		GroupId: gid,
		UserId:  uid,
	}
	lock.RLock()
	clients[wsKey] = ws
	lock.RUnlock()

	if err != nil {
		fmt.Println(err)
		delete(clients, wsKey) //删除map中的客户端
		return
	}
	defer ws.Close()
	for {
		msg := new(GroupSendMsg)
		//读取websocket发来的数据
		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Println("数据格式不正确", err)
			delete(clients, wsKey) //删除map中的客户端
			break
		}

		//发送消息
		if msg.Type == 1 {

			if len(msg.Content) == 0 {
				//跳过
				continue
			}

			//创建基础聊天消息模板
			chatMessage := GroupMessage{
				MessageData: msg.Content,
				UserId:      uid,
				GroupId:     gid,
				CreateTime:  utils.GetLocalDateTime(),
				ExpireTime:  utils.GetExpireDateTime(),
			}
			//如果消息为空

			//查询用户信息，获取用户头像与昵称
			user := model.User{}
			model.DB.Model(&user).Where("id=?", chatMessage.UserId).Find(&user)

			//创建用户聊天消息模板（往基础聊天消息模板里添加用户头像与用户昵称）
			userMessage := UserMessage{
				//Message中的数据
				MessageData: chatMessage.MessageData,
				User:        serializer.BuildUser(user),
				GroupId:     chatMessage.GroupId,
				CreateTime:  chatMessage.CreateTime,
			}

			//go groupMysqlSave(chatMessage)
			go groupRedisSave(chatMessage)
			//将聊天消息模板（用户发送的消息及用户头像昵称）插入广播频道，推送给其他websocket客户端
			groupBroadcast <- userMessage

			//获取历史消息
		} else if msg.Type == 2 {
			name := strconv.FormatInt(gid, 10)
			messages := cache.RedisClient.LRange(context.TODO(), name, 0, -1).Val()
			var groupMessages []model.GroupMessage
			for _, message := range messages {

				var groupMessage model.GroupMessage
				json.Unmarshal([]byte(message), &groupMessage)

				//过期删除
				t1, err1 := utils.TranslateTime(groupMessage.ExpireTime)
				if err1 != nil {
					fmt.Println(err1)
				}
				t2 := utils.GetLocalTime()
				if t2.After(t1) {
					cache.RedisClient.LRem(context.TODO(), name, 1, message)
					continue
				}

				groupMessages = append(groupMessages, groupMessage)
			}

			content, _ := json.Marshal(serializer.BuildGroupMessages(groupMessages))
			replyMsg := SingleMessage{
				FromID:     "0",
				FromUser:   "服务器",
				Code:       e.WebsocketSuccessMessage,
				Content:    string(content),
				CreateTime: utils.GetLocalDateTime(),
			}

			groupMessage, err1 := json.Marshal(replyMsg)
			if err1 != nil {
				fmt.Println("json格式转换失败", err)
			}

			err = ws.WriteMessage(websocket.TextMessage, groupMessage)
			if err1 != nil {
				fmt.Println("InsertOneMsg Err", err)
			}
		}
	}
}

// 广播推送消息
func pushMessages() {
	for {
		//读取通道中的消息
		msg := <-groupBroadcast

		//轮询现有的websocket客户端
		for key, client := range clients {

			//获取该客户端的群组id（GroupId）
			gid := key.GroupId

			//匹配客户端，判断该客户端的GroupId是否与该消息的GroupId一致，如果是，则将该消息投递给该客户端
			if msg.GroupId == gid && msg.User.ID != 0 && len(msg.MessageData) > 0 {
				//发送消息，含失败重试机制，重试次数:3
				for i := 0; i < 3; i++ {
					//发送消息到消费者客户端
					err := client.WriteJSON(msg)
					//如果发送成功
					if err == nil {
						//结束循环
						break
					}
					//如果到达重试次数，但仍未发送成功
					if i == 2 && err != nil {
						fmt.Println("InsertOneMsg Err", err)
						//客户端关闭
						client.Close()
						//删除map中的客户端
						delete(clients, key)
					}
				}
			}
		}
	}
}
func groupMysqlSave(replyMsg GroupMessage) {
	var msg model.GroupMessage
	msg.MessageData = replyMsg.MessageData
	msg.CreateTime = replyMsg.CreateTime
	msg.GroupId = replyMsg.GroupId
	msg.UserId = replyMsg.UserId
	msg.CreateTime = replyMsg.CreateTime
	err := model.DB.Save(&msg).Error
	if err != nil {
		fmt.Println("消息mysql存储失败！！！")
	}
}

// redis储存
func groupRedisSave(replyMsg GroupMessage) {
	var msg model.GroupMessage
	msg.MessageData = replyMsg.MessageData
	msg.CreateTime = replyMsg.CreateTime
	msg.GroupId = replyMsg.GroupId
	msg.UserId = replyMsg.UserId
	msg.ExpireTime = replyMsg.ExpireTime
	message, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json格式转换失败")
	}
	_, err = cache.RedisClient.RPush(context.TODO(), strconv.FormatInt(replyMsg.GroupId, 10), string(message)).Result()
	if err != nil {
		fmt.Println("消息mysql存储失败！！！")
	}

	//发送消息
	cache.GroupStreamMQ.SendMsg(context.Background(), &cache.Msg{
		Topic: "group",
		Body:  message,
	})
}
