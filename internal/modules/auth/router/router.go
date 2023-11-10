package router

import (
	"project_layout/internal/modules/auth/handler"

	"github.com/gin-gonic/gin"
)

func BindAuthRoutes(
	authGroup *gin.RouterGroup,
	h handler.HTTPHandler,
) {
	authGroup.POST("/signup", h.SignUp())
}
