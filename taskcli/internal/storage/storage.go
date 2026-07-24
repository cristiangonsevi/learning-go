package storage

import "taskcli/internal/model"

type Store interface {
	Save(t *model.Task) error
	List() ([]*model.Task, error)
}
