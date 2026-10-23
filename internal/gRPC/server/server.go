package grpc_server

import (
	"context"
	"io"
	"task-tracker/internal/gRPC/task_tracker_pb"
	"task-tracker/internal/job"
	"task-tracker/internal/service"
)

type Server struct {
	creator *service.Creator
	task_tracker_pb.UnimplementedTaskTrackerServiceServer
}

func NewServer(creator *service.Creator) *Server {
	return &Server{creator: creator}
}

func (s *Server) Add(ctx context.Context, in *task_tracker_pb.AddRequest) (*task_tracker_pb.AddResponse, error) {
	job, err := s.creator.NewJob(job.ADD)
	if err != nil {
		return nil, err
	}
	reader, err := job.Execute([]string{in.Name})
	if err != nil {
		return nil, err
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &task_tracker_pb.AddResponse{
		Result: string(result),
	}, nil
}
