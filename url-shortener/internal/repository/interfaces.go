package repository

import "context"

type UserStorageInterface interface {
	CreateUser(ctx context.Context) (string, error)
}
