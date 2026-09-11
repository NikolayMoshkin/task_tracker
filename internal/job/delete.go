package job

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/store"
)

type Delete struct {
	name    JobType
	handler store.DataHandler
}

func NewDeleteJob(h store.DataHandler) *Delete {
	return &Delete{DELETE, h}
}

func (j *Delete) Name() JobType {
	return j.name
}

func (j *Delete) Execute(params []string) (io.Reader, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", app_errors.InvalidArgs)
	}

	id, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, app_errors.InvalidArgs
	}

	if err := j.handler.Delete(id); err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task deleted. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Job = (*Delete)(nil)
