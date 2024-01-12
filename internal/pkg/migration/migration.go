package migration

import (
	"gorm.io/gorm"

	"project_layout/internal/models"
)

func Migrate(
	db *gorm.DB,
) error {
	return db.AutoMigrate(models.User{})
}
