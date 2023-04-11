package main

import (
	"chat/conf"
	v1 "chat/route"
	"chat/service/chat"
)

func main() {
	conf.Init()
	chat.InitChat()
	r := v1.NewRouter()
	_ = r.Run(conf.HttpPort)
}
