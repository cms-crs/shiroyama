package clients

import (
	"context"
	"fmt"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type UserClient struct {
	client userv1.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userv1.NewUserServiceClient(conn)
	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) GetUser(ctx context.Context, userID string) (*userv1.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &userv1.GetUserRequest{
		Id: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp, nil
}

func (c *UserClient) GetUsersByIds(ctx context.Context, userIDs []string) ([]*userv1.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &userv1.GetUsersByIdsRequest{
		Ids: userIDs,
	}

	resp, err := c.client.GetUsersByIds(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by ids: %w", err)
	}

	return resp.Users, nil
}
