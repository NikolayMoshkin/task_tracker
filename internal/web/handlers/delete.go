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

type Delete struct {
	cr *service.Creator
}

func NewDelete(cr *service.Creator) *Delete {
	return &Delete{cr}
}

func (a *Delete) Handle(w http.ResponseWriter, r *http.Request) error {
	j, err := a.cr.NewJob(job.DELETE)
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
