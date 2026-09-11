package job

import (
	"fmt"
	"io"
	"strings"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/helper"
	"task-tracker/internal/store"
	"time"
)

type List struct {
	name    JobType
	handler store.DataHandler
}

func NewListJob(h store.DataHandler) *List {
	return &List{LIST, h}
}

func (j *List) Name() JobType {
	return j.name
}

func (j *List) Execute(params []string) (io.Reader, error) {
	if len(params) > 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", app_errors.InvalidArgs)
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

var _ Job = (*List)(nil)
