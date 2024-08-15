package models

type Urls struct {
	Id    uint   `gorm:"primaryKey"`
	Alias string `gorm:"type:text;unique;not null"`
	Url   string `gorm:"not null"`
}
