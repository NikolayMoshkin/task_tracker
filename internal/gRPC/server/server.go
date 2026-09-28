package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"task-tracker/internal/gRPC/task_tracker_pb"
	"task-tracker/internal/job"
	"task-tracker/internal/service"
	"task-tracker/internal/store"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	task_tracker_pb.UnimplementedTaskTrackerServiceServer
}

func (s *server) Add(ctx context.Context, in *task_tracker_pb.AddRequest) (*task_tracker_pb.AddResponse, error) {
	name := in.Name

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

	job, err := cr.NewJob(job.ADD)
	if err != nil {
		log.Fatal(err)
	}
	reader, err := job.Execute([]string{name})
	if err != nil {
		log.Fatal(err)
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	return &task_tracker_pb.AddResponse{
		Result: string(result),
	}, nil

}

func main() {
	fmt.Println("server  starting...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051") // default port for gRPC
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	task_tracker_pb.RegisterTaskTrackerServiceServer(s, &server{})

	// рефлексия для клиентов по типу evans/grpcurl
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
