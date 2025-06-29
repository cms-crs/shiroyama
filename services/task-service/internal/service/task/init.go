package service

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"log/slog"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *taskv1.Task) (*taskv1.Task, error)
	GetTask(ctx context.Context, taskID string) (*taskv1.Task, error)
	UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
	MoveTask(ctx context.Context, req *taskv1.MoveTaskRequest) (*taskv1.Task, error)
	AssignUser(ctx context.Context, taskID, userID string) (*taskv1.Task, error)
	UnassignUser(ctx context.Context, taskID, userID string) (*taskv1.Task, error)
	GetTasksForLists(ctx context.Context, listIDs []string) ([]*taskv1.Task, error)
	GetTasksForUser(ctx context.Context, userID string) ([]*taskv1.Task, error)
}

type Service struct {
	log         *slog.Logger
	taskRepo    TaskRepository
	userClient  userv1.UserServiceClient
	boardClient boardv1.BoardServiceClient
}

func NewTaskService(
	log *slog.Logger,
	taskRepo TaskRepository,
	userClient userv1.UserServiceClient,
	boardClient boardv1.BoardServiceClient,
) *Service {
	return &Service{
		log:         log,
		taskRepo:    taskRepo,
		userClient:  userClient,
		boardClient: boardClient,
	}
}

func (s *Service) validateListExists(ctx context.Context, listID string) error {
	_, err := s.boardClient.UpdateList(ctx, &boardv1.UpdateListRequest{Id: listID})
	if err != nil {
		return fmt.Errorf("board not found: %w", err)
	}
	return nil
}

func (s *Service) validateBoardExists(ctx context.Context, boardID string) error {
	_, err := s.boardClient.GetBoard(ctx, &boardv1.GetBoardRequest{Id: boardID})
	if err != nil {
		return fmt.Errorf("board not found: %w", err)
	}
	return nil
}
