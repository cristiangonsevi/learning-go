package handler

import (
	"context"
	"net/http"
	"url-shortener/internal/model"

	"github.com/gin-gonic/gin"
)

type URLServiceInterface interface {
	CreateURL(ctx context.Context, params model.CreateURLParams)
}

type URLHandlerStruct struct {
	urlService URLServiceInterface
}

func NewURLHandler(urlService URLServiceInterface) *URLHandlerStruct {
	return &URLHandlerStruct{
		urlService,
	}
}

func (h *URLHandlerStruct) CreateURLHandler(c *gin.Context) {
	var req model.CreateURLParams

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Error with the payload",
			"errors":  transalateErrors(err),
		})
	}

	h.urlService.CreateURL(c.Request.Context(), req)

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Url created successfully",
		"data":    "",
	})
}
