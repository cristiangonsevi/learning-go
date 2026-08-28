package services

import (
	"context"
	"math/rand"
	"url-shortener/internal/model"
)

type URLRepositoryInterface interface {
	CreateURL(ctx context.Context, params model.CreateURLParams)
}

type URLServiceStruct struct {
	repo URLRepositoryInterface
}

func NewURLService(urlRepo URLRepositoryInterface) *URLServiceStruct {

	return &URLServiceStruct{
		repo: urlRepo,
	}
}

func (s *URLServiceStruct) CreateURL(ctx context.Context, params model.CreateURLParams) {
	_, _ = GenerateShortURL(params.URL)
}

func GenerateShortURL(url string) (string, string) {
	shortURL := make([]byte, 6)
	base62 := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for idx := range 6 {
		shortURL[idx] = base62[rand.Intn(len(base62))]
	}
	return url, string(shortURL)
}
