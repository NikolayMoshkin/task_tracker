package job

import (
	"fmt"
	"io"
	"strings"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/model"
	"task-tracker/internal/store"
)

type Add struct {
	name    JobType
	handler store.DataHandler
}

func NewAddJob(h store.DataHandler) *Add {
	return &Add{ADD, h}
}

func (j *Add) Name() JobType {
	return j.name
}

func (j *Add) Execute(params []string) (io.Reader, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", app_errors.InvalidArgs)
	}

	taskName := strings.TrimSpace(params[0])
	if taskName == "" {
		return nil, fmt.Errorf("Task name should not be empty%w", app_errors.InvalidArgs)
	}

	id, err := j.handler.Add(model.NewTask(params[0]))
	if err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task added. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Job = (*Add)(nil)
