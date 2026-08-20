package router

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")

	{
		auth.POST("/register")
	}

	return r
}
