package job

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/store"
)

type MarkInProgress struct {
	name    JobType
	handler store.DataHandler
}

func NewMarkInProgressJob(h store.DataHandler) *MarkInProgress {
	return &MarkInProgress{MARK_IN_PROGRESS, h}
}

func (j *MarkInProgress) Name() JobType {
	return j.name
}

func (j *MarkInProgress) Execute(params []string) (io.Reader, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", app_errors.InvalidArgs)
	}

	id, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	if err := j.handler.Update(id, "status", "in-progress"); err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task set 'in progress'. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Job = (*MarkInProgress)(nil)
