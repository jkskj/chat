package chat

// InitChat 聊天室模块初始化
func InitChat() {

	go pushMessages()

	go ClientsManager.Start()
}
