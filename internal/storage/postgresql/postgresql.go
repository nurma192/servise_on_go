package postgresql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"service_on_go/internal/storage/models"
)

type Storage struct {
	db *gorm.DB
}

func New() (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := gorm.Open("postgres", "user=postgres password=uk888888 dbname=urlshorter sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.AutoMigrate(&models.Urls{})

	return &Storage{db: db}, nil

}

func (s *Storage) SaveUrl(url string, alias string) (uint, error) {
	const op = "storage.postgres.New"

	newShortUrl := &models.Urls{Url: url, Alias: alias}

	if !s.db.NewRecord(newShortUrl) {
		return 0, fmt.Errorf("%s: this shorterUrl is already has in storage", op)
	}

	return newShortUrl.Id, s.db.Create(newShortUrl).Error
}

func (s *Storage) GetUrl() {

}

func (s *Storage) DeleteUrl() {

}
