package models

type Url struct {
	Id    uint   `gorm:"primaryKey"`
	Alias string `gorm:"not null"`
	Url   string `gorm:"not null"`
}
