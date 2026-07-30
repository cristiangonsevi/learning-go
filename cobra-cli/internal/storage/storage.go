package storage

import "cobra-cli/internal/model"

type Storage interface {
	List() ([]model.TaskModel, error)
	Create(model.TaskModel) error
}
