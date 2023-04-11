package model

type Join struct {
	ID  uint `gorm:"primary_key"`
	Gid uint `gorm:"not null"`
	Uid uint `gorm:"not null"`
}
