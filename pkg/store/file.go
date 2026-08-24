package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"task_tracker/pkg/helper"
	"task_tracker/pkg/model"
	"time"
)

type FileHandler struct {
	filename string
	file     *os.File
}

func NewFileHandler(filename ...string) (DataHandler, error) {
	fn := "tasks.json" // default filename
	if len(filename) > 0 {
		fn = filename[0]
	}

	return &FileHandler{
		filename: fn,
	}, nil
}

// returns id of a new task
func (h *FileHandler) Add(t *model.Task) (int, error) {
	f, close, err := h.getFile()
	if err != nil {
		return 0, err
	}
	defer close()

	newId, err := h.getNewId(f)
	if err != nil {
		return 0, err
	}

	t.Id = newId

	if err := h.addTask(f, t); err != nil {
		return 0, err
	}

	return newId, nil
}

func (h *FileHandler) Update(id int, field string, val string) error {
	f, close, err := h.getFile()
	if err != nil {
		return err
	}
	defer close()

	tasks, err := h.getTasks(f)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return errors.New("No tasks found")
	}

	var data []byte
	var updated bool
	for _, task := range tasks {
		if task.Id != id {
			continue
		}
		switch field {
		case "name":
			task.Name = val
		case "status":
			task.Status = val
		default:
			return errors.New("unknown task field: " + field)
		}

		now := time.Now()
		task.UpdatedAt = &now

		data, err = json.Marshal(tasks)
		if err != nil {
			return err
		}
		updated = true
		break
	}

	if updated == false {
		return fmt.Errorf("Task with id %d not found", id)
	}

	err = os.Truncate(f.Name(), 0)
	if err != nil {
		return err
	}
	f.Seek(0, 0)
	f.Write(data)

	return nil
}

func (h *FileHandler) Delete(id int) error {
	f, close, err := h.getFile()
	if err != nil {
		return err
	}
	defer close()

	tasks, err := h.getTasks(f)
	if err != nil {
		return err
	}

	var data []byte
	var updated bool
	for i, task := range tasks {
		if task.Id != id {
			continue
		}
		tasks = slices.Delete(tasks, i, i+1)
		updated = true

		data, err = json.Marshal(tasks)
		if err != nil {
			return err
		}
		break
	}

	if updated == false {
		return fmt.Errorf("Task with id %d not found", id)
	}

	err = os.Truncate(f.Name(), 0)
	if err != nil {
		log.Fatalf("Failed to truncate file: %v", err)
	}
	f.Seek(0, 0)
	f.Write(data)

	return nil
}

func (h *FileHandler) isExist() bool {
	_, err := os.Stat(h.filename)
	if err == nil {
		return true
	}
	return false
}

func (h *FileHandler) getFile() (*os.File, func(), error) {
	var f *os.File
	var err error

	if h.isExist() {
		f, err = os.OpenFile(h.filename, os.O_RDWR, 0644)
	} else {
		f, err = os.OpenFile(h.filename, os.O_CREATE|os.O_RDWR, 0644)
	}

	if err != nil {
		return nil, nil, err
	}

	closeFn := func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	return f, closeFn, nil
}

func (h *FileHandler) isEmpty(f *os.File) bool {
	finfo, _ := f.Stat()
	return finfo.Size() == 0
}

func (h *FileHandler) hasData(f *os.File) bool {
	if _, err := f.Seek(0, 0); err != nil {
		return false
	}
	defer f.Seek(0, 0)

	var b [1]byte
	if _, err := io.ReadFull(f, b[:]); err == nil && b[0] == '[' {
		return true
	}
	return false
}

// read all data and spit by tasks
func (h *FileHandler) getTasks(f *os.File) ([]*model.Task, error) {
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return []*model.Task{}, nil
	}

	var tasks []*model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (h *FileHandler) getNewId(f *os.File) (int, error) {
	newId := 1

	if !h.hasData(f) {
		return newId, nil
	}

	tasks, err := h.getTasks(f)
	if err != nil {
		return 0, err
	}

	maxId, err := helper.MaxKey(tasks)
	if err != nil {
		return newId, nil
	}

	return maxId + 1, nil
}

func (h *FileHandler) addTask(f *os.File, t *model.Task) error {
	tasks, err := h.getTasks(f)
	if err != nil {
		return err
	}

	tasks = append(tasks, t)

	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func (h *FileHandler) Select(field string, val string) ([]*model.Task, error) {
	f, close, err := h.getFile()
	if err != nil {
		return nil, err
	}
	defer close()

	tasks, err := h.getTasks(f)
	if err != nil {
		return nil, err
	}
	if val == "" || len(tasks) == 0 {
		return tasks, nil
	}

	var res []*model.Task
	for _, task := range tasks {
		switch field {
		case "name":
			if task.Name == val {
				res = append(res, task)
			}
		case "status":
			if task.Status == val {
				res = append(res, task)
			}
		default:
			return nil, errors.New("unknown task field: " + field)
		}
	}

	return res, nil
}

var _ DataHandler = (*FileHandler)(nil)
