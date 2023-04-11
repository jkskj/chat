package model

type GroupMessage struct {
	Id          int64  `gorm:"column:id;primary_key;auto_increment" json:"id"`
	GroupId     int64  `gorm:"column:group_id" json:"group_id"`
	UserId      int64  `gorm:"column:user_id" json:"user_id"`
	MessageData string `gorm:"column:message_data" json:"message_data"`
	CreateTime  string `gorm:"column:create_time" json:"create_time"`
	ExpireTime  string `gorm:"-" json:"expire_time"`
}

type SingleMessage struct {
	Id          int64  `gorm:"column:id;primary_key;auto_increment" json:"id"`
	ToId        int64  `gorm:"column:to_id" json:"to_id"`
	FromId      int64  `gorm:"column:from_id" json:"from_id"`
	MessageData string `gorm:"column:message_data" json:"message_data"`
	CreateTime  string `gorm:"column:create_time" json:"create_time"`
	ExpireTime  string `gorm:"-" json:"expire_time"`
}
