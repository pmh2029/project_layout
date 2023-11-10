package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB, logger *logrus.Logger) *gin.Engine {
	ginServer := gin.New()

	ginServer.Use(gin.Logger())
	ginServer.Use(gin.Recovery())
	ginServer.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowWebSockets:  true,
		AllowFiles:       true,
	}))

	ginServer.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	return ginServer
}
