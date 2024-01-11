package database

import (
	"fmt"
	"project_layout/configs"
	customLogger "project_layout/pkg/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

// NewDB initialize database
func NewDB(
	config configs.Config,
	logger *logrus.Logger,
) (*gorm.DB, error) {
	var dbConn *gorm.DB
	var err error

	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=prefer"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)

	dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: customLogger.NewGormLogger(logger),
	})
	if err != nil {
		return nil, err
	}

	err = Ping(dbConn)
	return dbConn, err
}

func CloseDB(
	logger *logrus.Logger,
	db *gorm.DB,
) {
	myDB, err := db.DB()
	if err != nil {
		logger.Errorf("Error while returning *sql.DB: %v", err)
	}

	logger.Info("Closing the DB connection pool")
	if err := myDB.Close(); err != nil {
		logger.Errorf("Error while closing the master DB connection pool: %v", err)
	}
}

func Ping(
	db *gorm.DB,
) error {
	myDB, err := db.DB()
	if err != nil {
		return err
	}

	return myDB.Ping()
}
