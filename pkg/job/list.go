package job

import (
	"fmt"
	"io"
	"strings"
	"task_tracker/pkg/helper"
	"task_tracker/pkg/store"
	"time"
)

type List struct {
	name    string
	handler store.DataHandler
}

func NewListJob(h store.DataHandler) *List {
	return &List{"list", h}
}

func (j *List) Name() string {
	return j.name
}

func (j *List) Execute(params []string) (io.Reader, error) {
	if len(params) > 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", InvalidArgs)
	}

	val := ""
	if len(params) == 1 {
		val = params[0]
	}

	tasks, err := j.handler.Select("status", val)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return strings.NewReader("No tasks found"), nil
	}

	output := ""
	for _, task := range tasks {
		output = output + fmt.Sprintf(`Id: %d
Name: %s,
Status: %s,
Created at: %s,
Updated at: %s
`,
			task.Id,
			task.Name,
			task.Status,
			task.CreatedAt.Format(time.DateTime),
			helper.FormatOptionalTime(task.UpdatedAt, time.DateTime),
		)
	}

	return strings.NewReader(output), nil
}

var _ Executable = (*List)(nil)
