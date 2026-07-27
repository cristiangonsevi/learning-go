package sqlite

import (
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	return &Store{db}, nil
}

func (store *Store) List() interface{} {
	return ""
}
