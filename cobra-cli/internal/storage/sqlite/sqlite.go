package sqlite

import (
	"cobra-cli/internal/model"
	"database/sql"
	"fmt"

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

func (store *Store) Find(id int) (*model.TaskModel, error) {

	var taskItem = model.TaskModel{}

	err := store.db.QueryRow("SELECT id, name, status, created_at FROM tasks WHERE id = ?", id).Scan(&taskItem.ID, &taskItem.Name, &taskItem.Status, &taskItem.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("No existe tarea con id = %v", id)
	}

	if err != nil {
		return nil, fmt.Errorf("Error obteniendo tarea\n")
	}

	return &taskItem, nil
}

func (store *Store) Update(id int, data model.TaskModel) error {

	taskItem, err := store.Find(id)

	if err != nil {
		return err
	}

	if data.Name != "" {
		taskItem.Name = data.Name
	}

	if data.Status != "" {
		taskItem.Status = data.Status
	}

	_, errUpdate := store.db.Exec("UPDATE tasks SET name = ?, status = ? WHERE id = ?", taskItem.Name, taskItem.Status, taskItem.ID)

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (store *Store) Delete(id int) error {

	_, err := store.Find(id)

	if err != nil {
		return err
	}

	_, err2 := store.db.Exec("DELETE FROM tasks WHERE id = ?", id)

	if err2 != nil {
		return err2
	}
	return nil
}

func (store *Store) Close() {
	store.db.Close()
}
