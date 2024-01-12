package database

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"project_layout/configs"
	customLogger "project_layout/pkg/logger"
)

// NewDB initialize database
func NewDB(
	config configs.Config,
	logger *logrus.Logger,
) (*gorm.DB, error) {
	var dbConn *gorm.DB
	var err error

	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=prefer"
	dsnMaster := fmt.Sprintf("%s:%s@tcp(%s:%s)/", config.DBUser, config.DBPass, config.DBHost, config.DBMasterPort)
	dsnSlave := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPass, config.DBHost, config.DBSlavePort, config.DBName)

	dbConn, err = gorm.Open(mysql.Open(dsnMaster))
	if err != nil {
		return nil, err
	}

	dsnMaster = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPass, config.DBHost, config.DBMasterPort, config.DBName)

	err = dbConn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DBName)).Error
	if err != nil {
		return nil, err
	}

	dbConn, err = gorm.Open(mysql.Open(dsnMaster), &gorm.Config{
		Logger: customLogger.NewGormLogger(customLogger.GormLogger{
			Logger:                    logger.WithField("service", "database"),
			LogLevel:                  gormLog.Info,
			IgnoreRecordNotFoundError: true,
			SlowThreshold:             200 * time.Millisecond,
			FileWithLineNumField:      "",
		}),
	})
	if err != nil {
		return nil, err
	}

	err = dbConn.Use(dbresolver.Register(dbresolver.Config{
		Replicas:          []gorm.Dialector{mysql.Open(dsnSlave)},
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: true,
	}))
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
