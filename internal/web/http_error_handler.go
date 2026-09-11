package web

import (
	"errors"
	"log/slog"
	"net/http"
	app_errors "task-tracker/internal/err"
)

type HttpErrorHandler struct {
	functionToWrap func(w http.ResponseWriter, r *http.Request) error
	logger         *slog.Logger
}

func NewHttpErrorHandler(f func(w http.ResponseWriter, r *http.Request) error, l *slog.Logger) *HttpErrorHandler {
	return &HttpErrorHandler{
		functionToWrap: f,
		logger:         l,
	}
}

func (eh *HttpErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := eh.functionToWrap(w, r); err != nil {
		if errors.Is(err, app_errors.ValidationError) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		eh.logger.Error("job execution error", "error", err)
	}
}
