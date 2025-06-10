package handler

import (
	"authservice/src/dto"
	"context"
	auth "github.com/cms-crs/protos/gen/go/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
}

type server struct {
	auth.UnimplementedAuthServiceServer
	logger      *slog.Logger
	authService AuthService
}

func RegisterServer(gRPC *grpc.Server, authService AuthService, logger *slog.Logger) {
	authServer := server{
		logger:      logger,
		authService: authService,
	}

	auth.RegisterAuthServiceServer(gRPC, &authServer)
}

func (s *server) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	const op = "grpc.Register"

	log := s.logger.With(
		slog.String("op", op),
	)

	registerRequest := dto.RegisterRequest{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	if err := registerRequest.Validate(); err != nil {
		log.Debug("validate create user failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// service call
	registerResponse, err := s.authService.Register(ctx, registerRequest)
	if err != nil {
		log.Debug("create user failed", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.RegisterResponse{
		UserId: uint64(registerResponse.UserId),
	}, nil
}

func (s *server) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	const op = "grpc.Login"

	log := s.logger.With(slog.String("op", op))

	loginRequest := dto.LoginRequest{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	if err := loginRequest.Validate(); err != nil {
		log.Debug("validate login failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// service call
	loginResponse, err := s.authService.Login(ctx, loginRequest)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.LoginResponse{
		AccessToken: loginResponse.Token,
	}, nil
}
