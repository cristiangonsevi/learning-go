package services

import (
	"context"
	"fmt"
	"time"
	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"url-shortener/internal/services/jwt"
	"uuid"

	"github.com/alexedwards/argon2id"
)

type UserInterface interface {
	CreateUser(ctx context.Context, user model.CreateUserRequest) (*model.UserAuthResponse, error)
	LoginUser(ctx context.Context, params model.LoginRequest) (*model.UserAuthResponse, error)
}

type UserService struct {
	repo repository.UserStorageInterface
}

func NewUserService(repo repository.UserStorageInterface) *UserService {
	return &UserService{
		repo,
	}
}

func (r *UserService) CreateUser(ctx context.Context, user model.CreateUserRequest) (*model.UserAuthResponse, error) {

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

	refreshToken, err := jwt.GenerateRefreshToken(payload, (time.Now().Add(360 * time.Hour).Unix()))

	if err != nil {
		return nil, err
	}

	_, err = r.repo.CreateUser(ctx, newUser)

	if err != nil {
		return nil, err
	}

	response := &model.UserAuthResponse{
		User: model.UserResponse{
			ID:    newUser.ID,
			Name:  newUser.Name,
			Email: newUser.Email,
		},
		JWT: model.TokenResponse{
			Token:        token,
			RefreshToken: refreshToken,
		},
	}

	return response, nil
}

func (r *UserService) LoginUser(ctx context.Context, params model.LoginRequest) (*model.UserAuthResponse, error) {

	user, err := r.repo.LoginUser(ctx, params)

	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(params.Password, user.Password)

	if err != nil {
		return nil, err
	}

	if !match {
		return nil, fmt.Errorf("Error: Bad credentials")
	}

	payload := map[string]interface{}{
		"uuid":  user.ID,
		"email": user.Email,
	}

	token, err := jwt.GenerateJWT(payload)

	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(payload, (time.Now().Add(360 * time.Hour).Unix()))

	if err != nil {
		return nil, err
	}

	response := &model.UserAuthResponse{
		User: model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		JWT: model.TokenResponse{
			Token:        token,
			RefreshToken: refreshToken,
		},
	}

	return response, nil
}
