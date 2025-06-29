package clients

import (
	"context"
	"fmt"
	"time"

	teamv1 "github.com/cms-crs/protos/gen/go/team_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClients struct {
	UserClient userv1.UserServiceClient
	TeamClient teamv1.TeamServiceClient
	userConn   *grpc.ClientConn
	teamConn   *grpc.ClientConn
}

type ClientConfig struct {
	UserServiceAddr string
	TeamServiceAddr string
	DialTimeout     time.Duration
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

	teamConn, err := grpc.DialContext(ctx, cfg.TeamServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		userConn.Close()
		return nil, fmt.Errorf("failed to connect to team service: %w", err)
	}

	return &ServiceClients{
		UserClient: userv1.NewUserServiceClient(userConn),
		TeamClient: teamv1.NewTeamServiceClient(teamConn),
		userConn:   userConn,
		teamConn:   teamConn,
	}, nil
}

func (c *ServiceClients) Close() error {
	var errUser, errTeam error

	if c.userConn != nil {
		errUser = c.userConn.Close()
	}

	if c.teamConn != nil {
		errTeam = c.teamConn.Close()
	}

	if errUser != nil {
		return fmt.Errorf("failed to close user service connection: %w", errUser)
	}

	if errTeam != nil {
		return fmt.Errorf("failed to close team service connection: %w", errTeam)
	}

	return nil
}
