package postgresql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"service_on_go/internal/storage"
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

	db.AutoMigrate(&models.Url{})

	return &Storage{db: db}, nil

}

func (s *Storage) SaveUrl(url string, alias string) (uint, error) {
	const op = "storage.postgres.New"

	isAliasExist := s.isExistByAlias(alias)
	if isAliasExist {
		return 0, fmt.Errorf("%s: this shorterUrl is already has in storage", op)
	}

	newShortUrl := &models.Url{Url: url, Alias: alias}
	if !s.db.NewRecord(newShortUrl) {
		return 0, fmt.Errorf("%s: DatabaseConnit save url with this parametres", op)
	}

	err := s.db.Create(newShortUrl).Error
	if err != nil {
		fmt.Println("some error when try to save the url to database")
	} else {
		fmt.Println("Url already saved")
	}

	fmt.Printf("url %s us saved", alias)

	return newShortUrl.Id, err
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.postgres.GetUrl"
	var url models.Url

	err := s.db.First(&url, &models.Url{Alias: alias}).Error

	if err != nil {
		return "", storage.ErrURLNotFound
	}

	return url.Url, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.postgres.DeleteUrl"

	err := s.db.Where("alias = ?", alias).Delete(&models.Url{}).Error

	if err != nil {
		return fmt.Errorf("%w: %s", err, op)
	}

	return nil
}
