package storage

import "taskcli/internal/model"

type Store interface {
	Save(t *model.Task) error
	List() ([]*model.Task, error)
	Delete(id *int) error
	Find(id *int) (*model.Task, error)
}
