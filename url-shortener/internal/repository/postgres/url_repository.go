package postgres

import (
	"context"
	"url-shortener/internal/model"
)

func (s *PostgresStorage) CreateURL(ctx context.Context, params model.URL) {
	s.db.NamedExecContext(ctx, `INSERT INTO url (id, url, short_url, user_id, title, description, is_active, expires_at) 
	VALUES (:id,:url,:short_url,:user_id,:title,:description,:is_active,:expires_at)`, &params)
}
