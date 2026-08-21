package services

import (
	"context"
	"url-shortener/internal/repository"
)

type UserInterface interface {
	CreateUser(ctx context.Context) (string, error)
}

type UserSevice struct {
	repo repository.UserStorageInterface
}

func NewUserService(repo repository.UserStorageInterface) *UserSevice {
	return &UserSevice{
		repo,
	}
}

func (r *UserSevice) CreateUser(ctx context.Context) (string, error) {
	r.repo.CreateUser(ctx)
	return "", nil
}
