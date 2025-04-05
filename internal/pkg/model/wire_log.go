package model

type WireLog struct {
	Import string `gorm:"not null"`
	Func   string `gorm:"not null"`
}
