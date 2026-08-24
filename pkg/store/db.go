package store

import (
	"task_tracker/pkg/model"
)

// database handler

type DBHandler struct {
}

func NewDBHandler() (DataHandler, error) {
	return &DBHandler{}, nil
}

// returns id of a new task
func (h *DBHandler) Add(t *model.Task) (int, error) {
	return 1, nil
}

func (h *DBHandler) Update(id int, field string, val string) error {
	return nil
}
func (h *DBHandler) Delete(id int) error {
	return nil
}

func (h *DBHandler) Select(field string, val string) ([]*model.Task, error) {
	var tasks []*model.Task
	return tasks, nil
}

var _ DataHandler = (*DBHandler)(nil)
