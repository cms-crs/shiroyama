package handler

import (
	"context"
	"github.com/cms-crs/protos/gen/go/user_service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"log/slog"
	"time"
	"userservice/internal/dto"
	"userservice/internal/kafka"
)

type UserService interface {
	CreateUser(
		context.Context,
		*dto.CreateUserRequest,
	) (*dto.CreateUserResponse, error)
	GetUser(
		context.Context,
		*dto.GetUserRequest,
	) (*dto.GetUserResponse, error)
	GetUsersByIds(
		context.Context,
		[]string,
	) (*dto.GetUsersByIdsResponse, error)
	UpdateUser(
		context.Context,
		*dto.UpdateUserRequest,
	) (*dto.UpdateUserResponse, error)
}

type GrpcHandler struct {
	userv1.UnimplementedUserServiceServer
	log           *slog.Logger
	userService   UserService
	kafkaProducer *kafka.Producer
}

func Register(gRPC *grpc.Server, log *slog.Logger, userService UserService, kafkaProducer *kafka.Producer) {
	userv1.RegisterUserServiceServer(gRPC, GrpcHandler{
		log:           log,
		userService:   userService,
		kafkaProducer: kafkaProducer,
	})
}

func (handler GrpcHandler) CreateUser(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.User, error) {
	const op = "gRPC.CreateUser"

	createUserRequest := &dto.CreateUserRequest{
		Email:    request.Email,
		Username: request.Username,
	}

	err := createUserRequest.Validate()
	if err != nil {
		handler.log.Error(op, err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	createUserResponse, err := handler.userService.CreateUser(ctx, createUserRequest)
	if err != nil {
		handler.log.Error(op, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userv1.User{
		Id:        createUserResponse.ID,
		Email:     createUserResponse.Email,
		Username:  createUserResponse.Username,
		CreatedAt: timestamppb.New(createUserResponse.CreatedAt),
		UpdatedAt: timestamppb.New(createUserResponse.UpdatedAt),
	}, nil
}

func (handler GrpcHandler) GetUser(ctx context.Context, request *userv1.GetUserRequest) (*userv1.User, error) {
	const op = "gRPC.GetUser"

	getUserRequest := &dto.GetUserRequest{ID: request.Id}

	err := getUserRequest.Validate()
	if err != nil {
		handler.log.Error(op, err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	getUserResponse, err := handler.userService.GetUser(ctx, getUserRequest)
	if err != nil {
		handler.log.Error(op, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userv1.User{
		Id:        getUserResponse.ID,
		Email:     getUserResponse.Email,
		Username:  getUserResponse.Username,
		CreatedAt: timestamppb.New(getUserResponse.CreatedAt),
		UpdatedAt: timestamppb.New(getUserResponse.UpdatedAt),
	}, nil
}

func (handler GrpcHandler) GetUsersByIds(ctx context.Context, request *userv1.GetUsersByIdsRequest) (*userv1.GetUsersByIdsResponse, error) {
	const op = "gRPC.GetUsersByIds"

	getUsersByIdsRequest := &dto.GetUsersByIdsRequest{IDs: request.Ids}

	err := getUsersByIdsRequest.Validate()
	if err != nil {
		handler.log.Error(op, err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	getUsersByIdsResponse, err := handler.userService.GetUsersByIds(ctx, getUsersByIdsRequest.IDs)
	if err != nil {
		handler.log.Error(op, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	var users []*userv1.User
	for _, user := range getUsersByIdsResponse.Users {
		users = append(users, &userv1.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &userv1.GetUsersByIdsResponse{
		Users: users,
	}, nil
}

func (handler GrpcHandler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.User, error) {
	const op = "gRPC.UpdateUser"

	updateUserRequest := &dto.UpdateUserRequest{
		ID:       req.Id,
		Username: req.Username,
	}

	updateUserResponse, err := handler.userService.UpdateUser(ctx, updateUserRequest)
	if err != nil {
		handler.log.Error(op, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userv1.User{
		Id:        updateUserResponse.ID,
		Username:  updateUserResponse.Username,
		Email:     updateUserResponse.Email,
		CreatedAt: timestamppb.New(updateUserResponse.CreatedAt),
		UpdatedAt: timestamppb.New(updateUserResponse.UpdatedAt),
	}, nil
}

func (handler GrpcHandler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*emptypb.Empty, error) {
	getUserRequest := &dto.GetUserRequest{ID: req.Id}
	user, err := handler.userService.GetUser(ctx, getUserRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	sagaID := uuid.New().String()
	event := kafka.Event{
		ID:        uuid.New().String(),
		Type:      kafka.UserDeletionRequested,
		UserID:    req.Id,
		Timestamp: time.Now(),
		SagaID:    sagaID,
		Data: map[string]interface{}{
			"user_email":    user.Email,
			"user_username": user.Username,
			"initiated_by":  "grpc_request",
		},
	}

	if err := handler.kafkaProducer.PublishEvent("user-deletion-saga", event); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to initiate user deletion: %v", err)
	}

	log.Printf("User deletion saga initiated for user %s (saga: %s)", req.Id, sagaID)
	return &emptypb.Empty{}, nil
}
