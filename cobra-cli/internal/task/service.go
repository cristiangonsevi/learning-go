package task

import (
	"cobra-cli/internal/model"
	"cobra-cli/internal/storage"
	"context"
)

type Service struct {
	store storage.Storage
}

func NewService(store storage.Storage) *Service {
	return &Service{store}
}

func (s *Service) List(ctx context.Context) ([]model.TaskModel, error) {
	tasks, err := s.store.List(ctx)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) CreateTask(ctx context.Context, newTask model.TaskModel) error {
	err := s.store.Create(ctx, newTask)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateTask(ctx context.Context, id int, data model.TaskModel) error {

	err := s.store.Update(ctx, id, data)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	err := s.store.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
