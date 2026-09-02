package repository

import (
	"context"
	"url-shortener/internal/model"
)

type UserStorageInterface interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	LoginUser(ctx context.Context, params model.LoginRequest) (*model.UserAuthResponse, error)
}
