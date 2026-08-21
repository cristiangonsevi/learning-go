package services

import (
	"context"
	"url-shortener/internal/repository"
)

type UserInterface interface {
	CreateUser(ctx context.Context) (string, error)
}

type UserService struct {
	repo repository.UserStorageInterface
}

func NewUserService(repo repository.UserStorageInterface) *UserService {
	return &UserService{
		repo,
	}
}

func (r *UserService) CreateUser(ctx context.Context) (string, error) {
	r.repo.CreateUser(ctx)
	return "", nil
}
