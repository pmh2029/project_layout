package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project_layout/internal/migration"
	"project_layout/internal/pkg/infrastructure"
	"project_layout/internal/pkg/mount"
	"project_layout/pkg/database"
	"project_layout/pkg/logger"
	"syscall"
	"time"
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

	ginServer := infrastructure.NewServer(db, logger)

	err = mount.MountAll(ginServer, db, logger)
	if err != nil {
		logger.Fatalln("Failed to mount the dependencies: ", err.Error())
		panic(err)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: ginServer,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	q := <-quit
	logger.Printf("Received signal '%v'. Shutting down...", q.String())

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err.Error())
	}

	logger.Println("Server exiting")
}
