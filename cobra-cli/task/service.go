package task

import (
	"cobra-cli/storage"
	"fmt"
)

type Service struct {
	store storage.Storage
}

func NewService(store storage.Storage) *Service {
	return &Service{store}
}

func (s *Service) List() interface{} {
	fmt.Println("Listar tareas")
	return s.store.List()
}
