package services

import (
	"context"
	"math/rand"
	"url-shortener/internal/model"
	"uuid"
)

type URLRepositoryInterface interface {
	CreateURL(ctx context.Context, params model.URL)
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
	shortUrl := GenerateShortURL(params.URL)

	url := model.URL{
		ID:          uuid.NewV7().String(),
		URL:         params.URL,
		ShortURL:    shortUrl,
		Title:       params.Title,
		Description: params.Description,
		ExpiresAt:   params.ExpiresAt,
		UserID:      params.UserID,
		IsActive:    params.IsActive,
	}

	s.repo.CreateURL(ctx, url)

}

func GenerateShortURL(url string) string {
	shortURL := make([]byte, 6)
	base62 := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for idx := range 6 {
		shortURL[idx] = base62[rand.Intn(len(base62))]
	}
	return string(shortURL)
}
