package task

import (
	"taskcli/internal/model"
	"taskcli/internal/storage"
)

type Service struct {
	store storage.Store
}

func NewService(store storage.Store) *Service {
	return &Service{store}
}

func (s *Service) CreateTask(t *model.Task) error {
	err := s.store.Save(t)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteTask(id *int) error {
	err := s.store.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FindTask(id *int) (*model.Task, error) {
	taskItem, err := s.store.Find(id)

	if err != nil {
		return nil, err
	}
	return taskItem, nil
}

func (s *Service) ListTasks() ([]*model.Task, error) {

	tasks, err := s.store.List()

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
