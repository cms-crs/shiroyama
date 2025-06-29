package handlers

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type UserHandler struct {
	grpcClients *clients.GRPCClients
	logger      logger.Logger
}

func NewUserHandler(grpcClients *clients.GRPCClients, log logger.Logger) *UserHandler {
	return &UserHandler{
		grpcClients: grpcClients,
		logger:      log,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.CreateUserRequest true "Create user request"
// @Success 201 {object} models.Response{data=models.User}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Creating user", "email", req.Email, "username", req.Username)

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)
	response, err := userClient.CreateUser(c.Request.Context(), &userv1.CreateUserRequest{
		Email:    req.Email,
		Username: req.Username,
	})

	if err != nil {
		h.logger.Error("Failed to create user", "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.AlreadyExists:
				utils.ErrorResponse(c, http.StatusConflict, "User already exists")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	user := &models.User{
		ID:        response.Id,
		Email:     response.Email,
		Username:  response.Username,
		CreatedAt: response.CreatedAt.AsTime(),
		UpdatedAt: response.UpdatedAt.AsTime(),
	}

	h.logger.Info("User created successfully", "user_id", user.ID)
	utils.CreatedResponse(c, user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} models.Response{data=models.User}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := utils.GetParamID(c, "id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	h.logger.Info("Getting user", "user_id", userID)

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)
	response, err := userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{
		Id: userID,
	})

	if err != nil {
		h.logger.Error("Failed to get user", "user_id", userID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "User not found")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user")
		return
	}

	user := &models.User{
		ID:        response.Id,
		Email:     response.Email,
		Username:  response.Username,
		CreatedAt: response.CreatedAt.AsTime(),
		UpdatedAt: response.UpdatedAt.AsTime(),
	}

	utils.SuccessResponse(c, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRequest true "Update user request"
// @Success 200 {object} models.Response{data=models.User}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := utils.GetParamID(c, "id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var req models.UpdateUserRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	currentUserID, exists := utils.GetUserIDFromContext(c)
	if !exists || currentUserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own profile")
		return
	}

	h.logger.Info("Updating user", "user_id", userID)

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)
	response, err := userClient.UpdateUser(c.Request.Context(), &userv1.UpdateUserRequest{
		Id:       userID,
		Username: req.Username,
	})

	if err != nil {
		h.logger.Error("Failed to update user", "user_id", userID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "User not found")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	user := &models.User{
		ID:        response.Id,
		Email:     response.Email,
		Username:  response.Username,
		CreatedAt: response.CreatedAt.AsTime(),
		UpdatedAt: response.UpdatedAt.AsTime(),
	}

	h.logger.Info("User updated successfully", "user_id", userID)
	utils.SuccessResponse(c, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user account
// @Tags users
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 204 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := utils.GetParamID(c, "id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	currentUserID, exists := utils.GetUserIDFromContext(c)
	if !exists || currentUserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only delete your own account")
		return
	}

	h.logger.Info("Deleting user", "user_id", userID)

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)
	_, err := userClient.DeleteUser(c.Request.Context(), &userv1.DeleteUserRequest{
		Id: userID,
	})

	if err != nil {
		h.logger.Error("Failed to delete user", "user_id", userID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "User not found")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	h.logger.Info("User deletion initiated successfully", "user_id", userID)
	utils.NoContentResponse(c)
}

// GetUsersByIDs godoc
// @Summary Get users by IDs
// @Description Get multiple users by their IDs
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.GetUsersByIDsRequest true "Get users by IDs request"
// @Success 200 {object} models.Response{data=models.GetUsersByIDsResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/users/batch [post]
func (h *UserHandler) GetUsersByIDs(c *gin.Context) {
	var req models.GetUsersByIDsRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Getting users by IDs", "count", len(req.IDs))

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)
	response, err := userClient.GetUsersByIds(c.Request.Context(), &userv1.GetUsersByIdsRequest{
		Ids: req.IDs,
	})

	if err != nil {
		h.logger.Error("Failed to get users by IDs", "count", len(req.IDs), "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get users")
		return
	}

	users := make([]models.User, 0, len(response.Users))
	for _, user := range response.Users {
		users = append(users, models.User{
			ID:        user.Id,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.AsTime(),
			UpdatedAt: user.UpdatedAt.AsTime(),
		})
	}

	apiResponse := &models.GetUsersByIDsResponse{
		Users: users,
	}

	utils.SuccessResponse(c, apiResponse)
}
