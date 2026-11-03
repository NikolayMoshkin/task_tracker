package web

import (
	"log/slog"
	"net/http"
	"net/http/pprof"
	"task-tracker/internal/service"
	http_handlers "task-tracker/internal/web/handlers"
)

type Router struct {
	serviceCreator *service.Creator
	logger         *slog.Logger
}

func NewRouter(cr *service.Creator, l *slog.Logger) *Router {
	return &Router{
		serviceCreator: cr,
		logger:         l,
	}
}

func (r *Router) Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", http_handlers.Health)
	mux.HandleFunc("GET /list", NewHttpErrorHandler(http_handlers.NewList(r.serviceCreator).Handle, r.logger).ServeHTTP)
	mux.HandleFunc("POST /add", NewHttpErrorHandler(http_handlers.NewAdd(r.serviceCreator).Handle, r.logger).ServeHTTP)
	mux.HandleFunc("POST /update", NewHttpErrorHandler(http_handlers.NewUpdate(r.serviceCreator).Handle, r.logger).ServeHTTP)
	mux.HandleFunc("DELETE /delete/{id}", NewHttpErrorHandler(http_handlers.NewDelete(r.serviceCreator).Handle, r.logger).ServeHTTP)
	mux.HandleFunc("POST /mark-done/{id}", NewHttpErrorHandler(http_handlers.NewMarkDone(r.serviceCreator).Handle, r.logger).ServeHTTP)
	mux.HandleFunc("POST /mark-in-progress/{id}", NewHttpErrorHandler(http_handlers.NewMarkInProgress(r.serviceCreator).Handle, r.logger).ServeHTTP)

	return mux
}

// profiler
func (r *Router) Pprof() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	return mux
}
