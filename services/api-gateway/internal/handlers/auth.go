package handlers

import (
	"context"
	"net/http"
	"time"

	"api-gateway/internal/clients"
	"api-gateway/internal/config"
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"

	authv1 "github.com/cms-crs/protos/gen/go/auth_service"
)

type AuthHandler struct {
	grpcClients *clients.GRPCClients
	config      *config.Config
	logger      logger.Logger
}

func NewAuthHandler(grpcClients *clients.GRPCClients, cfg *config.Config, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		grpcClients: grpcClients,
		config:      cfg,
		logger:      log,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, and username
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.Response{data=models.RegisterResponse}
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("User registration attempt", "email", req.Email, "username", req.Username)

	authClient := h.grpcClients.GetAuthClient()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	response, err := authClient.Register(ctx, &authv1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	})
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	h.logger.Info("User registered successfully", "user_id", response.UserId)
	utils.CreatedResponse(c, &models.RegisterResponse{
		UserID: response.UserId,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.Response{data=models.LoginResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("User login attempt", "email", req.Email)

	authClient := h.grpcClients.GetAuthClient()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	response, err := authClient.Login(ctx, &authv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Error("Failed to login user", "error", err, "email", req.Email)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	h.logger.Info("User logged in successfully", "email", req.Email)
	utils.SuccessResponse(c, &models.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    h.config.AccessTokenDuration,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.Response{data=models.RefreshTokenResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Token refresh attempt")

	authClient := h.grpcClients.GetAuthClient()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	response, err := authClient.RefreshToken(ctx, &authv1.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		h.logger.Error("Failed to refresh token", "error", err)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	h.logger.Info("Token refreshed successfully")
	utils.SuccessResponse(c, &models.RefreshTokenResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    h.config.AccessTokenDuration,
	})
}

// ValidateToken godoc
// @Summary Validate token
// @Description Validate access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ValidateTokenRequest true "Validate token request"
// @Success 200 {object} models.Response{data=models.ValidateTokenResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/auth/validate [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req models.ValidateTokenRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	authClient := h.grpcClients.GetAuthClient()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	response, err := authClient.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		Token: req.Token,
	})
	if err != nil {
		h.logger.Error("Failed to validate token", "error", err)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	utils.SuccessResponse(c, &models.ValidateTokenResponse{
		Valid: response.Valid,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LogoutRequest true "Logout request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.LogoutRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("User logout")

	authClient := h.grpcClients.GetAuthClient()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	_, err := authClient.Logout(ctx, &authv1.LogoutRequest{
		Token: req.Token,
	})
	if err != nil {
		h.logger.Error("Failed to logout user", "error", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	h.logger.Info("User logged out successfully")
	utils.SuccessResponse(c, gin.H{"message": "Logged out successfully"})
}
