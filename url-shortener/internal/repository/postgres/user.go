package postgres

import (
	"context"
	"url-shortener/internal/model"
)

func (r *PostgresStorage) CreateUser(ctx context.Context, user model.CreateUserRequest) (string, error) {
	r.db.ExecContext(ctx, "")
	return "", nil
}
