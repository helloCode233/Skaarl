package model

type Cache struct {
	Key   string `gorm:"not null;index;"`
	Value string `gorm:"not null;"`
}
