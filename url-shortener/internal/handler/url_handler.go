package handler

import (
	"net/http"
	"time"
	"url-shortener/internal/model"

	"github.com/gin-gonic/gin"
)

type URLServiceInterface interface {
	Create(CreateURLRequest) (model.URL, error)
}

type URLHandlerStruct struct {
	urlService URLServiceInterface
}

type CreateURLRequest struct {
	URL       string
	ExpiresAt time.Time
}

func NewURLService(urlService URLServiceInterface) *URLHandlerStruct {
	return &URLHandlerStruct{
		urlService,
	}
}

func (h *URLHandlerStruct) CreateURLHanbler(c *gin.Context) {
	var req CreateURLRequest

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Error with the payload",
			"errors":  transalateErrors(err),
		})
	}

	h.urlService.Create(req)

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Url created successfully",
		"data":    "",
	})
}
