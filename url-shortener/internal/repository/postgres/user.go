package postgres

import "context"

func (r *PostgresStorage) CreateUser(ctx context.Context) (string, error) {
	r.db.ExecContext(ctx, "")
	return "", nil
}
