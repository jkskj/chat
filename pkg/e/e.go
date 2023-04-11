package e

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	ErrorExistUser       = 102
	ErrorNotExistUser    = 103
	ErrorFailEncryption  = 106
	ErrorNotCompare      = 107
	ErrorExistFriendship = 108
	ErrorNotExistGroup   = 109
	ErrorNotExistJoin    = 110

	ErrorAuthCheckTokenFail    = 301 //token 错误
	ErrorAuthCheckTokenTimeout = 302 //token 过期
	ErrorAuthToken             = 303
	ErrorAuth                  = 304
	ErrorDatabase              = 401
	ErrorNotExistData          = 402
	ErrorNotAdmin              = 403

	ErrorHaveJoinGroup = 405
	ErrorWebsocket     = 406

	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004
	WebsocketOfflineReply   = 50005
	WebsocketLimit          = 50006
)

var MsgFlags = map[int]string{
	SUCCESS:       "成功",
	ERROR:         "失败!!!!",
	InvalidParams: "请求参数错误",

	ErrorAuthCheckTokenFail:    "Token鉴权失败",
	ErrorAuthCheckTokenTimeout: "Token已超时",
	ErrorAuthToken:             "Token生成失败",
	ErrorAuth:                  "Token错误",
	ErrorNotExistUser:          "用户不存在，请先注册",
	ErrorNotCompare:            "密码不匹配",
	ErrorDatabase:              "数据库操作出错,请重试",
	ErrorExistUser:             "用户已存在",
	ErrorFailEncryption:        "加密密码失败",
	ErrorNotExistData:          "数据不存在",
	ErrorWebsocket:             "websocket升级失败",

	ErrorNotAdmin:        "不是管理员",
	ErrorExistFriendship: "好友已存在",
	ErrorHaveJoinGroup:   "已加入群聊，不可重复加入",
	ErrorNotExistGroup:   "群聊不存在",
	ErrorNotExistJoin:    "未加入群聊",

	WebsocketSuccessMessage: "解析content内容信息",
	WebsocketSuccess:        "发送信息，请求历史纪录操作成功",
	WebsocketEnd:            "请求历史纪录，但没有更多记录了",
	WebsocketOnlineReply:    "针对回复信息在线应答成功",
	WebsocketOfflineReply:   "针对回复信息离线回答成功",
	WebsocketLimit:          "请求收到限制",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
