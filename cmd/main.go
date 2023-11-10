package main

import (
	"net/http"
	"project_layout/internal/migration"
	"project_layout/internal/mount"
	"project_layout/internal/server"
	"project_layout/pkg/database"
	"project_layout/pkg/logger"
)

func main() {
	logger := logger.NewCustomLogger()

	logger.Info("Init Database")
	db, err := database.NewDB(logger)
	if err != nil {
		logger.Fatalln("Failed to connect database.")
		panic(err)
	}
	logger.Info("Init Database Success")

	defer database.CloseDB(logger, db)

	logger.Info("Migrate Database")
	err = migration.Migrate(db)
	if err != nil {
		logger.Fatalln("Failed to migrate database.")
		panic(err)
	}
	logger.Info("Migrate Database Success")

	ginServer := server.NewServer(db, logger)

	err = mount.MountAll(ginServer, db, logger)
	if err != nil {
		logger.Fatalln("Failed to mount the dependencies: ", err.Error())
		panic(err)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: ginServer,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalln("An error happened while starting the HTTP server: ", err)
	}
}
