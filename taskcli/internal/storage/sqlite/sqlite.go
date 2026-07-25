package sqlite

import (
	"database/sql"
	"fmt"
	"taskcli/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) List() ([]*model.Task, error) {
	rows, err := s.db.Query("SELECT id, name FROM task")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*model.Task

	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil

}

func (s *Store) Save(t *model.Task) error {
	_, err := s.db.Exec("INSERT INTO task (name) VALUES (?)", t.Name)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Delete(id *int) error {
	_, err := s.db.Exec("DELETE FROM task WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Find(id *int) (*model.Task, error) {
	var taskIdem model.Task

	err := s.db.QueryRow("SELECT id, name FROM task WHERE id = ?", id).Scan(&taskIdem.ID, &taskIdem.Name)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Tarea con id %d no encontrado", *id)
	}

	if err != nil {
		return nil, err
	}

	return &taskIdem, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
