package services

import (
	"context"
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

func (s *URLServiceStruct) CreateURL(ctx context.Context, params model.CreateURLParams) {}
