package handler

import (
	"net/http"
	"url-shortener/internal/model"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	urlService string
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
}
