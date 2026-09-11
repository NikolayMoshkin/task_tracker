package web

import (
	"fmt"
	"io"
	"net/http"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/job"
	"task-tracker/internal/service"
)

type List struct {
	cr *service.Creator
}

func NewList(cr *service.Creator) *List {
	return &List{cr}
}

func (l *List) Handle(w http.ResponseWriter, r *http.Request) error {
	j, err := l.cr.NewJob(job.LIST)
	if err != nil {
		return err
	}

	queryParams := r.URL.Query()
	searchQuery := queryParams.Get("status")

	reader, err := j.Execute([]string{searchQuery})
	if err != nil {
		return fmt.Errorf("%w: %s", app_errors.ValidationError, "status must not be empty")
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	return nil
}
