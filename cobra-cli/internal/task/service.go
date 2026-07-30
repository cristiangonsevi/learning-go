package task

import (
	"cobra-cli/internal/model"
	"cobra-cli/internal/storage"
	"fmt"
	"log"
)

type Service struct {
	store storage.Storage
}

func NewService(store storage.Storage) *Service {
	return &Service{store}
}

func (s *Service) List() {
	fmt.Println("Listar tareas")
	tasks, err := s.store.List()

	if err != nil {
		log.Fatal("Error retrieving data: ", err)
	}
	log.Println("Tasks: ", tasks)
}

func (s *Service) CreateTask(newTask model.TaskModel) error {
	err := s.store.Create(newTask)
	if err != nil {
		return err
	}
	return nil
}
