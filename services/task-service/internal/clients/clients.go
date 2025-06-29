package clients

import (
	"context"
	"fmt"
	"time"

	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClients struct {
	UserClient  userv1.UserServiceClient
	BoardClient boardv1.BoardServiceClient
	userConn    *grpc.ClientConn
	boardConn   *grpc.ClientConn
}

type ClientConfig struct {
	UserServiceAddr  string
	BoardServiceAddr string
	DialTimeout      time.Duration
}

func NewServiceClients(cfg *ClientConfig) (*ServiceClients, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	userConn, err := grpc.DialContext(ctx, cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	boardConn, err := grpc.DialContext(ctx, cfg.BoardServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		userConn.Close()
		return nil, fmt.Errorf("failed to connect to board service: %w", err)
	}

	return &ServiceClients{
		UserClient:  userv1.NewUserServiceClient(userConn),
		BoardClient: boardv1.NewBoardServiceClient(boardConn),
		userConn:    userConn,
		boardConn:   boardConn,
	}, nil
}

func (c *ServiceClients) Close() error {
	var errUser, errBoard error

	if c.userConn != nil {
		errUser = c.userConn.Close()
	}

	if c.boardConn != nil {
		errBoard = c.boardConn.Close()
	}

	if errUser != nil {
		return fmt.Errorf("failed to close user service connection: %w", errUser)
	}

	if errBoard != nil {
		return fmt.Errorf("failed to close board service connection: %w", errBoard)
	}

	return nil
}
