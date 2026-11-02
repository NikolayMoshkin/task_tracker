package task_tracker_grpc_server

import (
	"context"
	"io"
	"log/slog"
	"strconv"
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
	respCh, errCh := s.executeAsync(job.ADD, []string{in.Name})

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

func (s *Server) List(ctx context.Context, in *task_tracker_pb.ListRequest) (*task_tracker_pb.ListResponse, error) {
	respCh, errCh := s.executeAsync(job.LIST, []string{in.Status})

	select {
	case val := <-respCh:
		return &task_tracker_pb.ListResponse{
			Tasks: string(val),
		}, nil
	case err := <-errCh:
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "List task failed: "+err.Error())
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	}
}

func (s *Server) Update(ctx context.Context, in *task_tracker_pb.UpdateRequest) (*task_tracker_pb.UpdateResponse, error) {
	id := strconv.Itoa(int(in.Id))
	name := in.Name

	respCh, errCh := s.executeAsync(job.UPDATE, []string{id, name})

	select {
	case val := <-respCh:
		return &task_tracker_pb.UpdateResponse{
			Result: string(val),
		}, nil
	case err := <-errCh:
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Update task failed: "+err.Error())
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	}
}

func (s *Server) Delete(ctx context.Context, in *task_tracker_pb.DeleteRequest) (*task_tracker_pb.DeleteResponse, error) {
	respCh, errCh := s.executeAsync(job.DELETE, []string{strconv.Itoa(int(in.Id))})
	select {
	case val := <-respCh:
		return &task_tracker_pb.DeleteResponse{
			Result: string(val),
		}, nil
	case err := <-errCh:
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Update task failed: "+err.Error())
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	}
}

func (s *Server) MarkDone(ctx context.Context, in *task_tracker_pb.MarkDoneRequest) (*task_tracker_pb.MarkDoneResponse, error) {
	respCh, errCh := s.executeAsync(job.MARK_DONE, []string{strconv.Itoa(int(in.Id))})
	select {
	case val := <-respCh:
		return &task_tracker_pb.MarkDoneResponse{
			Result: string(val),
		}, nil
	case err := <-errCh:
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Update task failed: "+err.Error())
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	}
}

func (s *Server) MarkInProgress(ctx context.Context, in *task_tracker_pb.MarkInProgressRequest) (*task_tracker_pb.MarkInProgressResponse, error) {
	respCh, errCh := s.executeAsync(job.MARK_IN_PROGRESS, []string{strconv.Itoa(int(in.Id))})
	select {
	case val := <-respCh:
		return &task_tracker_pb.MarkInProgressResponse{
			Result: string(val),
		}, nil
	case err := <-errCh:
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "Update task failed: "+err.Error())
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	}
}

func (s *Server) executeAsync(j job.JobType, params []string) (<-chan []byte, <-chan error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	go func() {
		defer close(respCh)

		job, err := s.creator.NewJob(j)
		if err != nil {
			errCh <- err
			return
		}

		reader, err := job.Execute(params)
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

	return respCh, errCh
}
