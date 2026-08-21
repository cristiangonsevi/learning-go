package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func New(connString string) (*PostgresStorage, error) {
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{
		db,
	}, nil
}
