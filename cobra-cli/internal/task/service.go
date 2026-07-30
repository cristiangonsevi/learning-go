package task

import (
	"cobra-cli/internal/model"
	"cobra-cli/internal/storage"
)

type Service struct {
	store storage.Storage
}

func NewService(store storage.Storage) *Service {
	return &Service{store}
}

func (s *Service) List() ([]model.TaskModel, error) {
	tasks, err := s.store.List()

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) CreateTask(newTask model.TaskModel) error {
	err := s.store.Create(newTask)
	if err != nil {
		return err
	}
	return nil
}
