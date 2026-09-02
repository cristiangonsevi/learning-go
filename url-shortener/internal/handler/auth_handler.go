package handler

import (
	"log"
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
			"message": "Bad request",
			"errors":  transalateErrors(err),
		})
		return
	}
	resp, err := h.userService.CreateUser(ctx, req)

	if err != nil {
		log.Println("Error creating user ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error creating user" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Usuario creado correctamente",
		"data":    resp,
	})
}

func (h *UserHandler) LoginUserHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Error parsing response",
			"errors":  transalateErrors(err),
		})
		return
	}
	h.userService.LoginUser(c.Request.Context(), req)
}
