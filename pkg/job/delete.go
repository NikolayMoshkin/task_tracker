package job

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"task_tracker/pkg/store"
)

type Delete struct {
	name    string
	handler store.DataHandler
}

func NewDeleteJob(h store.DataHandler) *Delete {
	return &Delete{"delete", h}
}

func (j *Delete) Name() string {
	return j.name
}

func (j *Delete) Execute(params []string) (io.Reader, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", InvalidArgs)
	}

	id, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	if err := j.handler.Delete(id); err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task deleted. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Executable = (*Delete)(nil)
