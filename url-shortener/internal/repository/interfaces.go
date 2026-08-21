package repository

import (
	"context"
	"url-shortener/internal/model"
)

type UserStorageInterface interface {
	CreateUser(ctx context.Context, user model.CreateUserRequest) (string, error)
}
