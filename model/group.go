package model

type Group struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	Creator   uint   `gorm:"not null"`  //创建者
	MemberNum int    `gorm:"default:1"` //成员数量
	Avatar    string `gorm:"default:''"`
}
