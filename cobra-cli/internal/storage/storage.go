package storage

import (
	"cobra-cli/internal/model"
	"context"
)

type Storage interface {
	List(context.Context) ([]model.TaskModel, error)
	Create(context.Context, model.TaskModel) error
	Update(context.Context, int, model.TaskModel) error
	Delete(context.Context, int) error
}
