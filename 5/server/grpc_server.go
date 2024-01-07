package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "lesson15_16/5/api"
)

type taskServer struct {
	fileStorage *FileStorage
	pb.UnimplementedTaskServiceServer
}

func NewTaskServer(fileStorage *FileStorage) *taskServer {
	return &taskServer{
		fileStorage: fileStorage,
	}
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	task := Task{
		TaskID:   req.TaskId,
		TaskName: req.TaskName,
	}

	err := s.fileStorage.Save(ctx, task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating task: %v", err)
	}

	return &pb.TaskResponse{
		TaskId:   task.TaskID,
		TaskName: task.TaskName,
	}, nil
}

func (s *taskServer) GetTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	task, err := s.fileStorage.Get(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting task: %v", err)
	}

	return &pb.TaskResponse{
		TaskId:   task.TaskID,
		TaskName: task.TaskName,
	}, nil
}

func StartServer() {
	fileStorage := NewFileStorage("tasks.json")
	taskServer := NewTaskServer(fileStorage)

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("gRPC server is running on port 50051")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
