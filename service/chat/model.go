package chat

import (
	"chat/serializer"
	"github.com/gorilla/websocket"
	"sync"
)

// Client 用户
type Client struct {
	ID       string
	SendID   string
	UserName string // 用户名
	Socket   *websocket.Conn
	Send     chan SingleMessage
}

// Manager 用户管理
type Manager struct {
	Clients   map[string]*Client
	Login     chan *Client
	UnLogin   chan *Client
	Broadcast chan *Broadcast
	Reply     chan *Client
}

var ClientsManager = Manager{
	Clients:   make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Login:     make(chan *Client),
	UnLogin:   make(chan *Client),
	Reply:     make(chan *Client),
	Broadcast: make(chan *Broadcast),
}

// SingleMessage 单聊回复的消息
type SingleMessage struct {
	FromID     string `json:"from_id"`
	FromUser   string `json:"from_user"`
	ToID       string `json:"to_id"`
	Code       int    `json:"code"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
	ExpireTime string `json:"expire_time"`
}

// Broadcast 广播类，包括广播内容和源用户
type Broadcast struct {
	User    *Client
	Message SingleMessage
	Type    int
}

// SingleSendMsg 单聊发送消息的类型
type SingleSendMsg struct {
	ToUid   string `json:"to_uid"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// 连接的客户端,把每个客户端都放进来。Key为WsKey结构体(GroupId int64, UserId int64)。Value为websocket连接
var clients = make(map[WsKey]*websocket.Conn)

// 广播通道，用于广播推送群聊用户发送的消息(带缓冲区，提高并发速率)
var groupBroadcast = make(chan GroupMessage, 10000)

var lock sync.RWMutex

// GroupSendMsg 群聊发送消息的类型
type GroupSendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// GroupMessage 群聊包装后的消息
type GroupMessage struct {
	GroupId     int64           ` json:"group_id"`
	MessageData string          `json:"message_data"`
	CreateTime  string          `json:"create_time"`
	User        serializer.User `json:"user"`
}

// WsKey 群聊key
type WsKey struct {
	GroupId int64
	UserId  int64
}
