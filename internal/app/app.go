package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	grpc_server "task-tracker/internal/gRPC/server"
	"task-tracker/internal/gRPC/task_tracker_pb"
	"task-tracker/internal/job"
	"task-tracker/internal/model"
	"task-tracker/internal/service"
	"task-tracker/internal/store"
	"task-tracker/internal/web"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	storeMode := os.Getenv("STORE")

	var handler store.DataHandler
	switch storeMode {
	case "file":
		handler, err = store.NewFileHandler()
		if err != nil {
			log.Fatal(err)
		}
	case "db":
		connStr := os.Getenv("DB_CONNECT")
		handler, err = store.NewDBHandler(connStr)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Invalid mode")
		os.Exit(1)
	}

	err = handler.Init()
	if err != nil {
		log.Fatal(err)
	}

	cr := service.NewCreator(handler)

	mode := os.Getenv("MODE")
	switch mode {
	case "http":
		runHttp(cr)
	case "cli":
		runCli(cr)
	case "grpc":
		runGrpc(cr)
	}
}

func runHttp(cr *service.Creator) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      web.NewRouter(cr, logger).Init(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Контекст, который отменяется по SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// Ждём сигнал завершения
	<-ctx.Done()
	stop()
	logger.Info("shutdown signal received")

	// Даём активным запросам время на завершение
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}

func runCli(cr *service.Creator) {
	var err error
	args := os.Args
	args = slices.Delete(args, 0, 1)

	if len(args) < 1 {
		log.Fatal("Invalid input")
	}

	command := model.NewCommand(args[0], args[1:])
	job, err := cr.NewJob(job.JobType(command.Name()))
	if err != nil {
		if errors.Is(err, service.InvalidJobError) { // test error wrapping
			log.Fatalf("Validation error: %s", err)
		}
		log.Fatal(err)
	}

	output, err := job.Execute(command.Params())
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, output)
}

func runGrpc(cr *service.Creator) {
	lis, err := net.Listen("tcp", "0.0.0.0:50051") // default port for gRPC
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	task_tracker_pb.RegisterTaskTrackerServiceServer(s, grpc_server.NewServer(cr))

	// рефлексия для клиентов по типу evans/grpcurl
	reflection.Register(s)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		fmt.Printf("gRPC server listening on %s\n", lis.Addr())
		serveErr <- s.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		s.GracefulStop()
	case err := <-serveErr:
		if err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}
}
