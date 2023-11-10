package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface{}

type UserRepository struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

func NewUserRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) UserRepositoryInterface {
	return &UserRepository{
		DB:     db,
		Logger: logger,
	}
}
