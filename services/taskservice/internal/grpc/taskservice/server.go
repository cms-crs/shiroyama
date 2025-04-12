package taskservicegrpc

import (
	"context"
	"github.com/ShiroyamaY/protos/gen/go/taskservice"
	"google.golang.org/grpc"
)

type serverAPI struct {
	taskservice.UnimplementedTaskServiceServer
}

func Register(gRPC *grpc.Server) {
	taskservice.RegisterTaskServiceServer(gRPC, new(serverAPI))
}

func (server serverAPI) CreateTask(context.Context, *taskservice.CreateTaskRequest) (*taskservice.Task, error) {
	panic("implement me")
}
func (server serverAPI) GetTask(context.Context, *taskservice.GetTaskRequest) (*taskservice.Task, error) {
	panic("implement me")
}
func (server serverAPI) UpdateTask(context.Context, *taskservice.UpdateTaskRequest) (*taskservice.Task, error) {
	panic("implement me")
}
func (server serverAPI) DeleteTask(context.Context, *taskservice.DeleteTaskRequest) (*taskservice.DeleteTaskResponse, error) {
	panic("implement me")
}
