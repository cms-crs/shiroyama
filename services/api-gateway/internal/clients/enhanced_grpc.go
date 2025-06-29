package clients

import (
	"api-gateway/internal/config"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	"context"
	"crypto/tls"
	"fmt"
	authv1 "github.com/cms-crs/protos/gen/go/auth_service"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	teamv1 "github.com/cms-crs/protos/gen/go/team_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

type EnhancedGRPCClients struct {
	connections map[string]*grpc.ClientConn
	config      *config.Config
	logger      logger.Logger

	AuthClient  authv1.AuthServiceClient
	UserClient  userv1.UserServiceClient
	TeamClient  teamv1.TeamServiceClient
	BoardClient boardv1.BoardServiceClient
	TaskClient  taskv1.TaskServiceClient
	//ActivityClient activityv1.ActivityServiceClient
}

type GRPCConnectionConfig struct {
	MaxRecvMsgSize   int
	MaxSendMsgSize   int
	ConnectTimeout   time.Duration
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration
	RetryPolicy      *RetryPolicy
	EnableTLS        bool
	TLSSkipVerify    bool
}

type RetryPolicy struct {
	MaxAttempts       uint
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
}

func NewEnhancedGRPCClients(cfg *config.Config, log logger.Logger) (*EnhancedGRPCClients, error) {
	clients := &EnhancedGRPCClients{
		connections: make(map[string]*grpc.ClientConn),
		config:      cfg,
		logger:      log,
	}

	connConfig := &GRPCConnectionConfig{
		MaxRecvMsgSize:   4 * 1024 * 1024,
		MaxSendMsgSize:   4 * 1024 * 1024,
		ConnectTimeout:   10 * time.Second,
		KeepAliveTime:    30 * time.Second,
		KeepAliveTimeout: 5 * time.Second,
		RetryPolicy: &RetryPolicy{
			MaxAttempts:       3,
			InitialBackoff:    100 * time.Millisecond,
			MaxBackoff:        5 * time.Second,
			BackoffMultiplier: 2.0,
		},
		EnableTLS:     cfg.Environment == "production",
		TLSSkipVerify: cfg.Environment != "production",
	}

	services := map[string]string{
		"auth":  cfg.Services.AuthService,
		"user":  cfg.Services.UserService,
		"team":  cfg.Services.TeamService,
		"board": cfg.Services.BoardService,
		"task":  cfg.Services.TaskService,
		//"comment":  cfg.Services.CommentService,
		//"activity": cfg.Services.ActivityService,
	}

	for name, addr := range services {
		conn, err := clients.createEnhancedConnection(addr, connConfig)
		if err != nil {
			clients.Close()
			return nil, fmt.Errorf("failed to connect to %s service at %s: %w", name, addr, err)
		}
		clients.connections[name] = conn
		log.Info("Connected to gRPC service", "service", name, "address", addr)
	}

	clients.AuthClient = authv1.NewAuthServiceClient(clients.connections["auth"])
	clients.UserClient = userv1.NewUserServiceClient(clients.connections["user"])
	clients.TeamClient = teamv1.NewTeamServiceClient(clients.connections["team"])
	clients.BoardClient = boardv1.NewBoardServiceClient(clients.connections["board"])
	clients.TaskClient = taskv1.NewTaskServiceClient(clients.connections["task"])
	//clients.ActivityClient = activityv1.NewActivityServiceClient(clients.connections["activity"])

	log.Info("All gRPC clients initialized successfully")
	return clients, nil
}

func (c *EnhancedGRPCClients) createEnhancedConnection(addr string, config *GRPCConnectionConfig) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(config.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(config.MaxSendMsgSize),
		),

		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                config.KeepAliveTime,
			Timeout:             config.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),

		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  config.RetryPolicy.InitialBackoff,
				Multiplier: config.RetryPolicy.BackoffMultiplier,
				Jitter:     0.2,
				MaxDelay:   config.RetryPolicy.MaxBackoff,
			},
			MinConnectTimeout: config.ConnectTimeout,
		}),

		grpc.WithUnaryInterceptor(c.unaryClientInterceptor),
		grpc.WithStreamInterceptor(c.streamClientInterceptor),
	}

	if config.EnableTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.TLSSkipVerify,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *EnhancedGRPCClients) unaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	err := utils.WithRetry(ctx, func() error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}, 3, 100*time.Millisecond)

	duration := time.Since(start)

	if err != nil {
		c.logger.Error("gRPC unary call failed",
			"method", method,
			"duration", duration,
			"error", err,
		)
	} else {
		c.logger.Debug("gRPC unary call completed",
			"method", method,
			"duration", duration,
		)
	}

	return err
}

func (c *EnhancedGRPCClients) streamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	start := time.Now()

	stream, err := streamer(ctx, desc, cc, method, opts...)

	duration := time.Since(start)

	if err != nil {
		c.logger.Error("gRPC stream call failed",
			"method", method,
			"duration", duration,
			"error", err,
		)
	} else {
		c.logger.Debug("gRPC stream call started",
			"method", method,
			"duration", duration,
		)
	}

	return stream, err
}

func (c *EnhancedGRPCClients) GetAuthClient() authv1.AuthServiceClient {
	return c.AuthClient
}

func (c *EnhancedGRPCClients) GetUserClient() userv1.UserServiceClient {
	return c.UserClient
}

func (c *EnhancedGRPCClients) GetTeamClient() teamv1.TeamServiceClient {
	return c.TeamClient
}

func (c *EnhancedGRPCClients) GetBoardClient() boardv1.BoardServiceClient {
	return c.BoardClient
}

func (c *EnhancedGRPCClients) GetTaskClient() taskv1.TaskServiceClient {
	return c.TaskClient
}

//func (c *EnhancedGRPCClients) GetActivityClient() activityv1.ActivityServiceClient {
//	return c.ActivityClient
//}

func (c *EnhancedGRPCClients) HealthCheck(ctx context.Context) map[string]utils.GRPCHealthStatus {
	status := make(map[string]utils.GRPCHealthStatus)

	for name, conn := range c.connections {
		status[name] = utils.CheckGRPCHealth(ctx, name, func(ctx context.Context) error {
			state := conn.GetState()
			if state.String() != "READY" && state.String() != "IDLE" {
				return fmt.Errorf("connection state: %s", state.String())
			}
			return nil
		})
	}

	return status
}

func (c *EnhancedGRPCClients) Close() {
	for name, conn := range c.connections {
		if conn != nil {
			if err := conn.Close(); err != nil {
				c.logger.Error("Failed to close gRPC connection", "service", name, "error", err)
			} else {
				c.logger.Info("Closed gRPC connection", "service", name)
			}
		}
	}
}

func (c *EnhancedGRPCClients) GetConnectionStats() map[string]interface{} {
	stats := make(map[string]interface{})

	for name, conn := range c.connections {
		stats[name] = map[string]interface{}{
			"state":  conn.GetState().String(),
			"target": conn.Target(),
		}
	}

	return stats
}

func (c *EnhancedGRPCClients) WaitForReady(ctx context.Context) error {
	for name, conn := range c.connections {
		if !conn.WaitForStateChange(ctx, conn.GetState()) {
			return fmt.Errorf("timeout waiting for %s service to be ready", name)
		}
		c.logger.Info("Service ready", "service", name)
	}
	return nil
}
