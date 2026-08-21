package handler

import (
	"net/http"
	"url-shortener/internal/model"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserInterface
}

func NewAuthHandler(userService services.UserInterface) *UserHandler {
	return &UserHandler{
		userService,
	}
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var ctx = c.Request.Context()
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	h.userService.CreateUser(ctx)
}
