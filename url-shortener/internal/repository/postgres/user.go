package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"url-shortener/internal/model"
)

func (r *PostgresStorage) CreateUser(ctx context.Context, user model.User) (string, error) {
	existUser, err := r.GetUserByEmail(ctx, user.Email)

	if err != nil {
		return "", err
	}

	if existUser != nil && existUser.Email == user.Email {
		return "", fmt.Errorf("Email already exist")
	}

	_, err = r.db.NamedExecContext(ctx, "INSERT INTO users (id, name, email, password) VALUES (:id, :name, :email, :password)", user)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (*model.UserResponse, error) {
	var user model.UserResponse
	err := s.db.GetContext(ctx, &user, "SELECT id, name, email FROM users WHERE email = $1", email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}
