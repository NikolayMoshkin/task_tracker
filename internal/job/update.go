package job

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/store"
)

type Update struct {
	name    JobType
	handler store.DataHandler
}

func NewUpdateJob(h store.DataHandler) *Update {
	return &Update{UPDATE, h}
}

func (j *Update) Name() JobType {
	return j.name
}

func (j *Update) Execute(params []string) (io.Reader, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("%w. Invalid amount of args", app_errors.InvalidArgs)
	}

	id, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	fmt.Println(params)

	if err := j.handler.Update(id, "name", params[1]); err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task updated. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Job = (*Update)(nil)
