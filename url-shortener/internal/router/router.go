package router

import (
	"url-shortener/internal/handler"

	"github.com/gin-gonic/gin"
)

func New(userHandler *handler.UserHandler, urlHandler *handler.URLHandlerStruct) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")

	{
		auth.POST("/register", userHandler.CreateUserHandler)
	}

	url := r.Group("/api/url")

	{
		url.POST("/", urlHandler.CreateURLHandler)
	}

	return r
}
