package clients

import (
	"context"
	"fmt"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc"
)

type UserClient struct {
	client userv1.UserServiceClient
}

func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		client: userv1.NewUserServiceClient(conn),
	}
}

func (uc *UserClient) GetUser(ctx context.Context, id string) (*userv1.User, error) {
	resp, err := uc.client.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return resp, nil
}
