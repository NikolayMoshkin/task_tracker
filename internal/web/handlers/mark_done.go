package web

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/job"
	"task-tracker/internal/service"
)

type MarkDone struct {
	cr *service.Creator
}

func NewMarkDone(cr *service.Creator) *MarkDone {
	return &MarkDone{cr}
}

func (a *MarkDone) Handle(w http.ResponseWriter, r *http.Request) error {
	j, err := a.cr.NewJob(job.MARK_DONE)
	if err != nil {
		return err
	}

	id := r.PathValue("id")
	if id == "" {
		return fmt.Errorf("Id is not set%w", app_errors.ValidationError)
	}

	reader, err := j.Execute([]string{id})
	if err != nil {
		if errors.Is(err, app_errors.InvalidArgs) {
			return fmt.Errorf("%s%w", err.Error(), app_errors.ValidationError)
		}

		return err
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	return nil
}
