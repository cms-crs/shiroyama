package router

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"api-gateway/internal/clients"
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
)

func New(cfg *config.Config, grpcClients *clients.GRPCClients, log logger.Logger) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logging(log))
	r.Use(middleware.CORS(cfg))

	authHandler := handlers.NewAuthHandler(grpcClients, cfg, log)
	userHandler := handlers.NewUserHandler(grpcClients, log)
	teamHandler := handlers.NewTeamHandler(grpcClients, log)
	boardHandler := handlers.NewBoardHandler(grpcClients, log)
	taskHandler := handlers.NewTaskHandler(grpcClients, log)

	r.GET("/health", healthCheck)
	r.GET("/ready", readinessCheck)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/validate", authHandler.ValidateToken)
	}

	protected := v1.Group("")
	protected.Use(middleware.Auth(cfg))

	users := protected.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/batch", userHandler.GetUsersByIDs)

		users.GET("/:id/teams", teamHandler.GetUserTeams)
		users.GET("/:id/boards", boardHandler.GetUserBoards)
	}

	teams := protected.Group("/teams")
	{
		teams.POST("", teamHandler.CreateTeam)
		teams.GET("/:team_id", teamHandler.GetTeam)
		teams.PUT("/:team_id", teamHandler.UpdateTeam)
		teams.DELETE("/:team_id", teamHandler.DeleteTeam)

		teams.POST("/:team_id/members", teamHandler.AddUserToTeam)
		teams.DELETE("/:team_id/members/:user_id", teamHandler.RemoveUserFromTeam)
		teams.PUT("/:team_id/members/:user_id/role", teamHandler.UpdateUserRole)
		teams.GET("/:team_id/members", teamHandler.GetTeamMembers)

		teams.GET("/:team_id/boards", boardHandler.GetTeamBoards)
	}

	boards := protected.Group("/boards")
	{
		boards.POST("", boardHandler.CreateBoard)
		boards.GET("/:id", boardHandler.GetBoard)
		boards.PUT("/:id", boardHandler.UpdateBoard)
		boards.DELETE("/:id", boardHandler.DeleteBoard)

		boards.POST("/:id/lists", boardHandler.CreateList)
		boards.PUT("/:id/lists/reorder", boardHandler.ReorderLists)
	}

	lists := protected.Group("/lists")
	{
		lists.PUT("/:id", boardHandler.UpdateList)
		lists.DELETE("/:id", boardHandler.DeleteList)
	}

	tasks := protected.Group("/tasks")
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}

	return r
}

func healthCheck(c *gin.Context) {
	utils.SuccessResponse(c, models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"api-gateway": "healthy",
		},
	})
}

func readinessCheck(c *gin.Context) {
	utils.SuccessResponse(c, models.HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Services: map[string]string{
			"api-gateway": "ready",
		},
	})
}
