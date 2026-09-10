package store

import "task-tracker/pkg/model"

type DataHandler interface {
	Init() error
	Add(t *model.Task) (int, error)
	Update(id int, field string, val string) error
	Delete(id int) error
	Select(field string, val string) ([]model.Task, error)
}
