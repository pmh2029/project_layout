package repository

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"project_layout/internal/models"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context) error
}

type UserRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserRepository(
	db *gorm.DB,
	logger *logrus.Logger,
) UserRepositoryInterface {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context) error {
	err := r.db.WithContext(ctx).Create(&models.User{
		Email:    strconv.Itoa(rand.Intn(9999999999999)),
		Username: strconv.Itoa(rand.Intn(9999999999999)),
	}).Error
	if err != nil {
		r.logger.Errorf("Create user error: %v", err)

		return err
	}
	return nil
}
