package job

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"task_tracker/pkg/store"
)

type MarkDone struct {
	name    string
	handler store.DataHandler
}

func NewMarkDoneJob(h store.DataHandler) *MarkDone {
	return &MarkDone{"mark-done", h}
}

func (j *MarkDone) Name() string {
	return j.name
}

func (j *MarkDone) Execute(params []string) (io.Reader, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("%w. Invalid amount of args", InvalidArgs)
	}

	id, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	if err := j.handler.Update(id, "status", "done"); err != nil {
		return nil, err
	}

	output := fmt.Sprintf("Task set 'done'. Id: %d\n", id)

	return strings.NewReader(output), nil
}

var _ Executable = (*MarkDone)(nil)
