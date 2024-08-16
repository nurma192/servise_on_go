package postgresql

import (
	"service_on_go/internal/storage/models"
)

func (s *Storage) isExistByAlias(alias string) bool {
	var url models.Url

	err := s.db.First(&url, models.Url{Alias: alias}).Error

	return err == nil && url.Id > 0
}
