package sqlite

import (
	"cobra-cli/internal/model"
	"context"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sqlx.DB
}

func New(dbPath string) (*Store, error) {

	db, err := sqlx.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	return &Store{db}, nil
}

func (store *Store) List(ctx context.Context) ([]model.TaskModel, error) {
	var taskItems []model.TaskModel
	err := store.db.SelectContext(ctx, &taskItems, "SELECT id, name, status, created_at FROM tasks")

	if err != nil {
		return nil, err
	}

	return taskItems, nil
}

func (store *Store) Create(ctx context.Context, newTask model.TaskModel) error {
	_, err := store.db.ExecContext(ctx, "INSERT INTO tasks (name, status, created_at) VALUES (?, ?, ?)", newTask.Name, newTask.Status, newTask.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) Find(ctx context.Context, id int) (*model.TaskModel, error) {

	var taskItem = model.TaskModel{}

	err := store.db.GetContext(ctx, &taskItem, "SELECT id, name, status, created_at FROM tasks WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	return &taskItem, nil
}

func (store *Store) Update(ctx context.Context, id int, data model.TaskModel) error {

	taskItem, err := store.Find(ctx, id)

	if err != nil {
		return err
	}

	if data.Name != "" {
		taskItem.Name = data.Name
	}

	if data.Status != "" {
		taskItem.Status = data.Status
	}

	_, errUpdate := store.db.ExecContext(ctx, "UPDATE tasks SET name = ?, status = ? WHERE id = ?", taskItem.Name, taskItem.Status, taskItem.ID)

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (store *Store) Delete(ctx context.Context, id int) error {

	_, err := store.Find(ctx, id)

	if err != nil {
		return err
	}

	_, err2 := store.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)

	if err2 != nil {
		return err2
	}
	return nil
}

func (store *Store) Close() {
	store.db.Close()
}
