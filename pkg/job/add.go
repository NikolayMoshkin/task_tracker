package job

import (
	"fmt"
	"io"
	"strings"
	"task_tracker/pkg/model"
	"task_tracker/pkg/store"
)

type Add struct {
	name    string
	handler store.DataHandler
}

func NewAddJob(h store.DataHandler) *Add {
	return &Add{"add", h}
}

func (j *Add) Name() string {
	return j.name
}

func (j *Add) Execute(params []string) (io.Reader, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", InvalidArgs)
	}

	id, err := j.handler.Add(model.NewTask(params[0]))
	if err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task added. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Executable = (*Add)(nil)
