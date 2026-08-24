package job

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"task_tracker/pkg/store"
)

type MarkInProgress struct {
	name    string
	handler store.DataHandler
}

func NewMarkInProgressJob(h store.DataHandler) *MarkInProgress {
	return &MarkInProgress{"mark-in-progress", h}
}

func (j *MarkInProgress) Name() string {
	return j.name
}

func (j *MarkInProgress) Execute(params []string) (io.Reader, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", InvalidArgs)
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

var _ Executable = (*MarkInProgress)(nil)
