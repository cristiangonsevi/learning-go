package sqlite

import (
	"cobra-cli/internal/model"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
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

func (store *Store) List() ([]model.TaskModel, error) {
	var taskItems []model.TaskModel
	rows, err := store.db.Query("SELECT id, name, status, created_at FROM tasks")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var taskItem model.TaskModel
		err := rows.Scan(&taskItem.ID, &taskItem.Name, &taskItem.Status, &taskItem.CreatedAt)
		taskItems = append(taskItems, taskItem)

		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return taskItems, nil
}

func (store *Store) Create(newTask model.TaskModel) error {
	_, err := store.db.Exec("INSERT INTO tasks (name, status, created_at) VALUES (?, ?, ?)", newTask.Name, newTask.Status, newTask.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
