package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	app_errors "task-tracker/internal/err"
	"task-tracker/internal/job"
	"task-tracker/internal/service"
)

type Update struct {
	cr *service.Creator
}

type updateRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewUpdate(cr *service.Creator) *Update {
	return &Update{cr}
}

func (a *Update) Handle(w http.ResponseWriter, r *http.Request) error {
	j, err := a.cr.NewJob(job.UPDATE)
	if err != nil {
		return err
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid json body%w", app_errors.ValidationError)
	}

	reader, err := j.Execute([]string{strconv.Itoa(req.Id), req.Name})
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
