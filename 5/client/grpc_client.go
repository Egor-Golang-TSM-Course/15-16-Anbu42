package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "lesson15_16/5/api"
)

func StartClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	createTaskResponse, err := client.CreateTask(context.Background(), &pb.TaskRequest{
		TaskId:   1,
		TaskName: "Task 1",
	})
	if err != nil {
		log.Fatalf("CreateTask failed: %v", err)
	}
	fmt.Printf("CreateTask response: %+v\n", createTaskResponse)

	getTaskResponse, err := client.GetTask(context.Background(), &pb.TaskRequest{})
	if err != nil {
		log.Fatalf("GetTask failed: %v", err)
	}
	fmt.Printf("GetTask response: %+v\n", getTaskResponse)

	time.Sleep(3 * time.Second)
}
