package job

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"task_tracker/pkg/store"
)

type Update struct {
	name    string
	handler store.DataHandler
}

func NewUpdateJob(h store.DataHandler) *Update {
	return &Update{"update", h}
}

func (j *Update) Name() string {
	return j.name
}

func (j *Update) Execute(params []string) (io.Reader, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("%w. Invalid amount of args", InvalidArgs)
	}

	id, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	if err := j.handler.Update(id, "name", params[1]); err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task updated. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Executable = (*Update)(nil)
