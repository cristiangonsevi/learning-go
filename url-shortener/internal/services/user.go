package services

import (
	"context"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"url-shortener/internal/services/jwt"
	"uuid"

	"github.com/alexedwards/argon2id"
)

type UserInterface interface {
	CreateUser(ctx context.Context, user model.CreateUserRequest) (*model.UserRegisterResponse, error)
}

type UserService struct {
	repo repository.UserStorageInterface
}

func NewUserService(repo repository.UserStorageInterface) *UserService {
	return &UserService{
		repo,
	}
}

func (r *UserService) CreateUser(ctx context.Context, user model.CreateUserRequest) (*model.UserRegisterResponse, error) {

	hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)

	if err != nil {
		return nil, err
	}

	newUser := model.User{
		ID:           uuid.NewV7().String(),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: hashedPassword,
	}

	payload := map[string]interface{}{
		"uuid":  newUser.ID,
		"email": newUser.Email,
	}

	token, err := jwt.GenerateJWT(payload)

	if err != nil {
		return nil, err
	}

	_, err = r.repo.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	response := &model.UserRegisterResponse{
		User: model.UserResponse{
			ID:    newUser.ID,
			Name:  newUser.Name,
			Email: newUser.Email,
		},
		Token:        token,
		RefreshToken: token,
	}

	return response, nil
}
