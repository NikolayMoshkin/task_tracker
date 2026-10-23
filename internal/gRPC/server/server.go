package task_tracker_grpc_server

import (
	"context"
	"io"
	"log/slog"
	"task-tracker/internal/gRPC/task_tracker_pb"
	"task-tracker/internal/job"
	"task-tracker/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	creator *service.Creator
	task_tracker_pb.UnimplementedTaskTrackerServiceServer
}

func NewServer(creator *service.Creator) *Server {
	return &Server{creator: creator}
}

func (s *Server) Add(ctx context.Context, in *task_tracker_pb.AddRequest) (*task_tracker_pb.AddResponse, error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	go func() {
		defer close(respCh)

		job, err := s.creator.NewJob(job.ADD)
		if err != nil {
			errCh <- err
			return
		}
		reader, err := job.Execute([]string{in.Name})
		if err != nil {
			errCh <- err
			return
		}

		result, err := io.ReadAll(reader)
		if err != nil {
			errCh <- err
			return
		}

		respCh <- result
	}()

	select {
	case val := <-respCh:
		return &task_tracker_pb.AddResponse{
			Result: string(val),
		}, nil
	case err := <-errCh:
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Add task failed: "+err.Error())
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	}
}
