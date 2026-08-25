package store

import "task_tracker/pkg/model"

type DataHandler interface {
	Add(t *model.Task) (int, error)
	Update(id int, field string, val string) error
	Delete(id int) error
	Select(field string, val string) ([]model.Task, error)
}
