package router

import (
	"url-shortener/internal/handler"

	"github.com/gin-gonic/gin"
)

func New(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")

	{
		auth.POST("/register", userHandler.CreateUserHandler)
	}

	return r
}
