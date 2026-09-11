package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/job"
	"task-tracker/internal/service"
)

type Add struct {
	cr *service.Creator
}

type addRequest struct {
	Name string `json:"name"`
}

func NewAdd(cr *service.Creator) *Add {
	return &Add{cr}
}

func (a *Add) Handle(w http.ResponseWriter, r *http.Request) error {
	j, err := a.cr.NewJob(job.ADD)
	if err != nil {
		return err
	}

	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid json body: %w", err)
	}

	reader, err := j.Execute([]string{req.Name})
	if err != nil {
		if errors.Is(err, app_errors.InvalidArgs) {
			return fmt.Errorf("%s%w", err.Error(), app_errors.ValidationError)
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	return nil
}
