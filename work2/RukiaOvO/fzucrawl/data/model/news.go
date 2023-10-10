package model

type News struct {
	Title  string `gorm:"not null"`
	Author string `gorm:"not null"`
	Text   string `gorm:"not null"`
	Date   string `gorm:"not null"`
	Nums   string `gorm:"not null"`
}
