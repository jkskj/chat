package model

type Friend struct {
	ID        uint  `gorm:"primary_key"`
	UserOneID uint  `gorm:"not null"`
	UserTwoID uint  `gorm:"not null"`
	IsPass    int64 `gorm:"default:0"`
}
