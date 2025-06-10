package handler

import (
	"authservice/src/dto"
	"context"
	auth "github.com/cms-crs/protos/go/gen/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type server struct {
	auth.UnimplementedAuthServiceServer
	logger *slog.Logger
}

func RegisterServer(gRPC *grpc.Server, logger *slog.Logger) {
	authServer := server{
		logger: logger,
	}

	auth.RegisterAuthServiceServer(gRPC, &authServer)
}

func (s *server) Register(_ context.Context, in *auth.RegisterRequest) (*auth.AuthResponse, error) {
	const op = "grpc.Register"

	log := s.logger.With(
		slog.String("op", op),
	)

	createUserDto := dto.CreateUserRequest{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	if err := createUserDto.Validate(); err != nil {
		log.Debug("validate create user failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &auth.AuthResponse{
		AccessToken:  "",
		RefreshToken: "",
		UserId:       0,
		ExpiresIn:    0,
	}, nil
}

//
//func (s *server) Login(ctx context.Context, in *auth.LoginRequest) (*auth.AuthResponse, error) {
//	panic("implement me")
//}
