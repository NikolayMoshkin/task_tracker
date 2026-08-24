package service

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"task_tracker/pkg/job"
	"task_tracker/pkg/model"
	"task_tracker/pkg/store"
)

var InvalidJobError = errors.New("Invalid job")

type Creator struct {
	jobs map[string]job.Executable
}

func NewCreator(h store.DataHandler) *Creator {
	cr := &Creator{
		jobs: make(map[string]job.Executable),
	}

	updateJob := job.NewUpdateJob(h)
	addJob := job.NewAddJob(h)
	deleteJob := job.NewDeleteJob(h)
	markInProgressJob := job.NewMarkInProgressJob(h)
	markDoneJob := job.NewMarkDoneJob(h)
	listJob := job.NewListJob(h)

	cr.jobs[updateJob.Name()] = updateJob
	cr.jobs[addJob.Name()] = addJob
	cr.jobs[deleteJob.Name()] = deleteJob
	cr.jobs[markInProgressJob.Name()] = markInProgressJob
	cr.jobs[markDoneJob.Name()] = markDoneJob
	cr.jobs[listJob.Name()] = listJob

	return cr
}

func (cr *Creator) NewJob(cmd *model.Command) (job.Executable, error) {
	if job, ok := cr.jobs[cmd.Name()]; ok {
		return job, nil
	}

	keys := strings.Join(slices.Collect(maps.Keys(cr.jobs)), ",")

	return nil, fmt.Errorf("%w \"%s\". Available jobs: %s", InvalidJobError, cmd.Name(), keys)
}
