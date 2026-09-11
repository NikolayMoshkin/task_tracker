package service

import (
	"errors"
	"fmt"
	"strings"
	"task-tracker/internal/job"
	"task-tracker/internal/store"
)

var InvalidJobError = errors.New("Invalid job")

type Creator struct {
	jobs map[job.JobType]job.Job
}

func NewCreator(h store.DataHandler) *Creator {
	cr := &Creator{
		jobs: make(map[job.JobType]job.Job),
	}

	// method set. Интерфейс Executable имплементирован только для указателей на джобы, а не на сами джобы
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

func (cr *Creator) NewJob(jobName job.JobType) (job.Job, error) {
	if job, ok := cr.jobs[jobName]; ok {
		return job, nil
	}

	avJobs := []string{}
	for _, job := range cr.jobs {
		avJobs = append(avJobs, string(job.Name()))
	}

	keys := strings.Join(avJobs, ",")

	return nil, fmt.Errorf("%w \"%s\". Available jobs: %s", InvalidJobError, jobName, keys)
}
