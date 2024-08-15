package postgresql

import (
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
		return nil, err
	}

	db.AutoMigrate(&models.Urls{})

	return &Storage{db: db}, nil

}
