package handler

import (
	"context"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"userservice/internal/dto"
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
	DeleteUser(
		context.Context,
		string,
	) error
}

type apiServer struct {
	userv1.UnimplementedUserServiceServer
	log         *slog.Logger
	userService UserService
}

func Register(gRPC *grpc.Server, log *slog.Logger, userService UserService) {
	userv1.RegisterUserServiceServer(gRPC, &apiServer{
		log:         log,
		userService: userService,
	})
}

func (apiServer *apiServer) CreateUser(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.User, error) {
	const op = "gRPC.CreateUser"

	createUserRequest := dto.NewCreateUserRequest(request)

	err := createUserRequest.Validate()
	if err != nil {
		apiServer.log.Error(op, err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	createUserResponse, err := apiServer.userService.CreateUser(ctx, createUserRequest)
	if err != nil {
		apiServer.log.Error(op, err.Error())
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

func (apiServer *apiServer) GetUser(ctx context.Context, request *userv1.GetUserRequest) (*userv1.User, error) {
	const op = "gRPC.GetUser"

	getUserRequest := &dto.GetUserRequest{ID: request.Id}

	err := getUserRequest.Validate()
	if err != nil {
		apiServer.log.Error(op, err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	getUserResponse, err := apiServer.userService.GetUser(ctx, getUserRequest)
	if err != nil {
		apiServer.log.Error(op, err.Error())
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

func (apiServer *apiServer) GetUsersByIds(ctx context.Context, request *userv1.GetUsersByIdsRequest) (*userv1.GetUsersByIdsResponse, error) {
	const op = "gRPC.GetUsersByIds"

	getUsersByIdsRequest := &dto.GetUsersByIdsRequest{IDs: request.Ids}

	err := getUsersByIdsRequest.Validate()
	if err != nil {
		apiServer.log.Error(op, err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	getUsersByIdsResponse, err := apiServer.userService.GetUsersByIds(ctx, getUsersByIdsRequest.IDs)
	if err != nil {
		apiServer.log.Error(op, err.Error())
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

func (apiServer *apiServer) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.User, error) {
	const op = "gRPC.UpdateUser"

	updateUserRequest := &dto.UpdateUserRequest{
		ID:       req.Id,
		Username: req.Username,
	}

	updateUserResponse, err := apiServer.userService.UpdateUser(ctx, updateUserRequest)
	if err != nil {
		apiServer.log.Error(op, err.Error())
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

func (apiServer *apiServer) DeleteUser(ctx context.Context, request *userv1.DeleteUserRequest) (*emptypb.Empty, error) {
	const op = "gRPC.DeleteUser"

	err := apiServer.userService.DeleteUser(ctx, request.Id)
	if err != nil {
		apiServer.log.Error(op, err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
