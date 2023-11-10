package container

import (
	authHandler "project_layout/internal/modules/auth/handler"
	authRouter "project_layout/internal/modules/auth/router"
	"project_layout/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

type HandlerContainer struct {
	AuthHandler *authHandler.HTTPHandler
}

func NewHandlerContainer(
	validator *validator.JsonSchemaValidator,
	graphql graphql.Schema,
	logger *logrus.Logger,
) HandlerContainer {
	authContainer := authHandler.NewAuthHandler(graphql, validator, logger)

	return HandlerContainer{
		AuthHandler: &authContainer,
	}
}

func SetupHandler(gin *gin.Engine, handler *HandlerContainer) {
	// versioning
	v1 := gin.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authRouter.BindAuthRoutes(authGroup, *handler.AuthHandler)
}
