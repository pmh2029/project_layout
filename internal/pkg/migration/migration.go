package migration

import (
	"project_layout/internal/models"

	"gorm.io/gorm"
)

func Migrate(
	db *gorm.DB,
) error {
	return db.AutoMigrate(models.User{})
}
