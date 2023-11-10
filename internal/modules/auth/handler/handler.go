package handler

import (
	"project_layout/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

type HTTPHandler interface {
	SignUp() gin.HandlerFunc
}

type AuthHandler struct {
	graphql   graphql.Schema
	validator *validator.JsonSchemaValidator
	logger    *logrus.Logger
}

func NewAuthHandler(
	graphql graphql.Schema,
	validator *validator.JsonSchemaValidator,
	logger *logrus.Logger,
) HTTPHandler {
	return &AuthHandler{
		graphql:   graphql,
		validator: validator,
		logger:    logger,
	}
}
