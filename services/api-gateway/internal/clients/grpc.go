package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"api-gateway/internal/config"

	authv1 "github.com/cms-crs/protos/gen/go/auth_service"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	teamv1 "github.com/cms-crs/protos/gen/go/team_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
)

type GRPCClients struct {
	connections map[string]*grpc.ClientConn
	AuthClient  authv1.AuthServiceClient
	UserClient  userv1.UserServiceClient
	TeamClient  teamv1.TeamServiceClient
	BoardClient boardv1.BoardServiceClient
	TaskClient  taskv1.TaskServiceClient
	//ActivityClient activityv1.ActivityServiceClient
}

func NewGRPCClients(cfg *config.Config) (*GRPCClients, error) {
	clients := &GRPCClients{
		connections: make(map[string]*grpc.ClientConn),
	}

	services := map[string]string{
		"auth":  cfg.Services.AuthService,
		"user":  cfg.Services.UserService,
		"team":  cfg.Services.TeamService,
		"board": cfg.Services.BoardService,
		"task":  cfg.Services.TaskService,
		//"activity": cfg.Services.ActivityService,
	}

	for name, addr := range services {
		conn, err := createGRPCConnection(addr)
		if err != nil {
			clients.Close()
			return nil, fmt.Errorf("failed to create GRPC connection to %s %s: %w", name, addr, err)
		}
		clients.connections[name] = conn
	}

	clients.AuthClient = authv1.NewAuthServiceClient(clients.connections["auth"])
	clients.UserClient = userv1.NewUserServiceClient(clients.connections["user"])
	clients.TeamClient = teamv1.NewTeamServiceClient(clients.connections["team"])
	clients.BoardClient = boardv1.NewBoardServiceClient(clients.connections["board"])
	clients.TaskClient = taskv1.NewTaskServiceClient(clients.connections["task"])
	//clients.CommentClient = commentv1.NewCommentServiceClient(clients.connections["comment"])
	//clients.ActivityClient = activityv1.NewActivityServiceClient(clients.connections["activity"])

	return clients, nil
}

func createGRPCConnection(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *GRPCClients) Close() {
	for _, conn := range c.connections {
		if conn != nil {
			conn.Close()
		}
	}
}

func (c *GRPCClients) GetAuthClient() authv1.AuthServiceClient {
	return c.AuthClient
}

func (c *GRPCClients) GetUserClient() userv1.UserServiceClient {
	return c.UserClient
}

func (c *GRPCClients) GetTeamClient() teamv1.TeamServiceClient {
	return c.TeamClient
}

func (c *GRPCClients) GetBoardClient() boardv1.BoardServiceClient {
	return c.BoardClient
}

func (c *GRPCClients) GetTaskClient() taskv1.TaskServiceClient {
	return c.TaskClient
}

//func (c *GRPCClients) GetActivityClient() activityv1.ActivityServiceClient {
//	return c.ActivityClient
//}

func (c *GRPCClients) HealthCheck(ctx context.Context) map[string]bool {
	status := make(map[string]bool)

	for name, conn := range c.connections {
		status[name] = conn.GetState().String() == "READY"
	}

	return status
}
