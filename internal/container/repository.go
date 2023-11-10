package container

import (
	userRepo "project_layout/internal/modules/auth/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	UserRepo userRepo.UserRepositoryInterface
}

func NewRepositoryContainer(
	db *gorm.DB,
	logger *logrus.Logger,
) RepositoryContainer {
	return RepositoryContainer{
		UserRepo: userRepo.NewUserRepository(db, logger),
	}
}
