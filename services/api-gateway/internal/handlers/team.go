package handlers

import (
	teamv1 "github.com/cms-crs/protos/gen/go/team_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"net/http"

	"api-gateway/internal/clients"
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TeamHandler struct {
	grpcClients *clients.GRPCClients
	logger      logger.Logger
}

func NewTeamHandler(grpcClients *clients.GRPCClients, log logger.Logger) *TeamHandler {
	return &TeamHandler{
		grpcClients: grpcClients,
		logger:      log,
	}
}

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new team
// @Tags teams
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.CreateTeamRequest true "Create team request"
// @Success 201 {object} models.Response{data=models.Team}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req models.CreateTeamRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	req.CreatedBy = userID

	h.logger.Info("Creating team", "name", req.Name, "created_by", req.CreatedBy)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	response, err := teamClient.CreateTeam(c.Request.Context(), &teamv1.CreateTeamRequest{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
	})

	if err != nil {
		h.logger.Error("Failed to create team", "error", err)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create team")
		return
	}

	team := &models.Team{
		ID:          response.Id,
		Name:        response.Name,
		Description: response.Description,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
	}

	h.logger.Info("Team created successfully", "team_id", team.ID)
	utils.CreatedResponse(c, team)
}

// GetTeam godoc
// @Summary Get team by ID
// @Description Get team information by ID
// @Tags teams
// @Produce json
// @Security Bearer
// @Param id path string true "Team ID"
// @Success 200 {object} models.Response{data=models.Team}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{id} [get]
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamID := utils.GetParamID(c, "id")
	if teamID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID is required")
		return
	}

	h.logger.Info("Getting team", "team_id", teamID)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	response, err := teamClient.GetTeam(c.Request.Context(), &teamv1.GetTeamRequest{
		Id: teamID,
	})

	if err != nil {
		h.logger.Error("Failed to get team", "team_id", teamID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Team not found")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get team")
		return
	}

	team := &models.Team{
		ID:          response.Id,
		Name:        response.Name,
		Description: response.Description,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
	}

	utils.SuccessResponse(c, team)
}

// UpdateTeam godoc
// @Summary Update team
// @Description Update team information
// @Tags teams
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Team ID"
// @Param request body models.UpdateTeamRequest true "Update team request"
// @Success 200 {object} models.Response{data=models.Team}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{id} [put]
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	teamID := utils.GetParamID(c, "id")
	if teamID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID is required")
		return
	}

	var req models.UpdateTeamRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Updating team", "team_id", teamID)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	response, err := teamClient.UpdateTeam(c.Request.Context(), &teamv1.UpdateTeamRequest{
		Id:          teamID,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		h.logger.Error("Failed to update team", "team_id", teamID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Team not found")
				return
			case codes.PermissionDenied:
				utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update team")
		return
	}

	team := &models.Team{
		ID:          response.Id,
		Name:        response.Name,
		Description: response.Description,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
	}

	h.logger.Info("Team updated successfully", "team_id", teamID)
	utils.SuccessResponse(c, team)
}

// DeleteTeam godoc
// @Summary Delete team
// @Description Delete team
// @Tags teams
// @Produce json
// @Security Bearer
// @Param id path string true "Team ID"
// @Success 204 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{id} [delete]
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	teamID := utils.GetParamID(c, "id")
	if teamID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID is required")
		return
	}

	h.logger.Info("Deleting team", "team_id", teamID)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	_, err := teamClient.DeleteTeam(c.Request.Context(), &teamv1.DeleteTeamRequest{
		Id: teamID,
	})

	if err != nil {
		h.logger.Error("Failed to delete team", "team_id", teamID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Team not found")
				return
			case codes.PermissionDenied:
				utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete team")
		return
	}

	h.logger.Info("Team deleted successfully", "team_id", teamID)
	utils.NoContentResponse(c)
}

// GetUserTeams godoc
// @Summary Get user teams
// @Description Get all teams for a user
// @Tags teams
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} models.Response{data=models.GetUserTeamsResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/users/{id}/teams [get]
func (h *TeamHandler) GetUserTeams(c *gin.Context) {
	userID := utils.GetParamID(c, "id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	h.logger.Info("Getting user teams", "user_id", userID)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	response, err := teamClient.GetUserTeams(c.Request.Context(), &teamv1.GetUserTeamsRequest{
		UserId: userID,
	})

	if err != nil {
		h.logger.Error("Failed to get user teams", "user_id", userID, "error", err)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user teams")
		return
	}

	teams := make([]models.TeamWithRole, 0, len(response.Teams))
	for _, teamWithRole := range response.Teams {
		teams = append(teams, models.TeamWithRole{
			Team: models.Team{
				ID:          teamWithRole.Team.Id,
				Name:        teamWithRole.Team.Name,
				Description: teamWithRole.Team.Description,
				CreatedAt:   teamWithRole.Team.CreatedAt.AsTime(),
				UpdatedAt:   teamWithRole.Team.UpdatedAt.AsTime(),
			},
			Role: teamWithRole.Role,
		})
	}

	apiResponse := &models.GetUserTeamsResponse{
		Teams: teams,
	}

	utils.SuccessResponse(c, apiResponse)
}

// AddUserToTeam godoc
// @Summary Add user to team
// @Description Add a user to a team with specific role
// @Tags teams
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Team ID"
// @Param request body models.AddUserToTeamRequest true "Add user to team request"
// @Success 201 {object} models.Response{data=models.TeamMember}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{id}/members [post]
func (h *TeamHandler) AddUserToTeam(c *gin.Context) {
	teamID := utils.GetParamID(c, "id")
	if teamID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID is required")
		return
	}

	var req models.AddUserToTeamRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Adding user to team", "team_id", teamID, "user_id", req.UserID, "role", req.Role)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	response, err := teamClient.AddUserToTeam(c.Request.Context(), &teamv1.AddUserToTeamRequest{
		TeamId: teamID,
		UserId: req.UserID,
		Role:   req.Role,
	})

	if err != nil {
		h.logger.Error("Failed to add user to team", "team_id", teamID, "user_id", req.UserID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				if st.Message() == "user already in team" {
					utils.ErrorResponse(c, http.StatusConflict, "User is already a member of this team")
					return
				}
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Team or user not found")
				return
			case codes.PermissionDenied:
				utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add user to team")
		return
	}

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)
	userResponse, err := userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{
		Id: req.UserID,
	})

	member := &models.TeamMember{
		TeamID: response.TeamId,
		UserID: response.UserId,
		Role:   response.Role,
	}

	if err == nil && userResponse != nil {
		member.User = models.User{
			ID:        userResponse.Id,
			Email:     userResponse.Email,
			Username:  userResponse.Username,
			CreatedAt: userResponse.CreatedAt.AsTime(),
			UpdatedAt: userResponse.UpdatedAt.AsTime(),
		}
	} else {
		h.logger.Warn("Failed to get user info for team member", "user_id", req.UserID, "error", err)
		// Возвращаем member без полной информации о пользователе
		member.User = models.User{
			ID: req.UserID,
		}
	}

	h.logger.Info("User added to team successfully", "team_id", teamID, "user_id", req.UserID)
	utils.CreatedResponse(c, member)
}

// RemoveUserFromTeam godoc
// @Summary Remove user from team
// @Description Remove a user from a team
// @Tags teams
// @Produce json
// @Security Bearer
// @Param team_id path string true "Team ID"
// @Param user_id path string true "User ID"
// @Success 204 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{team_id}/members/{user_id} [delete]
func (h *TeamHandler) RemoveUserFromTeam(c *gin.Context) {
	teamID := utils.GetParamID(c, "team_id")
	userID := utils.GetParamID(c, "user_id")

	if teamID == "" || userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID and User ID are required")
		return
	}

	h.logger.Info("Removing user from team", "team_id", teamID, "user_id", userID)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	_, err := teamClient.RemoveUserFromTeam(c.Request.Context(), &teamv1.RemoveUserFromTeamRequest{
		TeamId: teamID,
		UserId: userID,
	})

	if err != nil {
		h.logger.Error("Failed to remove user from team", "team_id", teamID, "user_id", userID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Team, user, or membership not found")
				return
			case codes.PermissionDenied:
				utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove user from team")
		return
	}

	h.logger.Info("User removed from team successfully", "team_id", teamID, "user_id", userID)
	utils.NoContentResponse(c)
}

// UpdateUserRole godoc
// @Summary Update user role in team
// @Description Update user's role in a team
// @Tags teams
// @Accept json
// @Produce json
// @Security Bearer
// @Param team_id path string true "Team ID"
// @Param user_id path string true "User ID"
// @Param request body models.UpdateUserRoleRequest true "Update user role request"
// @Success 200 {object} models.Response{data=models.TeamMember}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{team_id}/members/{user_id}/role [put]
func (h *TeamHandler) UpdateUserRole(c *gin.Context) {
	teamID := utils.GetParamID(c, "team_id")
	userID := utils.GetParamID(c, "user_id")

	if teamID == "" || userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID and User ID are required")
		return
	}

	var req models.UpdateUserRoleRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Updating user role in team", "team_id", teamID, "user_id", userID, "role", req.Role)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	response, err := teamClient.UpdateUserRole(c.Request.Context(), &teamv1.UpdateUserRoleRequest{
		TeamId: teamID,
		UserId: userID,
		Role:   req.Role,
	})

	if err != nil {
		h.logger.Error("Failed to update user role", "team_id", teamID, "user_id", userID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Team, user, or membership not found")
				return
			case codes.PermissionDenied:
				utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)
	userResponse, err := userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{
		Id: userID,
	})

	member := &models.TeamMember{
		TeamID: response.TeamId,
		UserID: response.UserId,
		Role:   response.Role,
	}

	if err == nil && userResponse != nil {
		member.User = models.User{
			ID:        userResponse.Id,
			Email:     userResponse.Email,
			Username:  userResponse.Username,
			CreatedAt: userResponse.CreatedAt.AsTime(),
			UpdatedAt: userResponse.UpdatedAt.AsTime(),
		}
	} else {
		h.logger.Warn("Failed to get user info for team member", "user_id", userID, "error", err)
		member.User = models.User{
			ID: userID,
		}
	}

	h.logger.Info("User role updated successfully", "team_id", teamID, "user_id", userID)
	utils.SuccessResponse(c, member)
}

// GetTeamMembers godoc
// @Summary Get team members
// @Description Get all members of a team
// @Tags teams
// @Produce json
// @Security Bearer
// @Param id path string true "Team ID"
// @Success 200 {object} models.Response{data=models.GetTeamMembersResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{id}/members [get]
func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	teamID := utils.GetParamID(c, "id")
	if teamID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID is required")
		return
	}

	h.logger.Info("Getting team members", "team_id", teamID)

	teamClient := h.grpcClients.GetTeamClient().(teamv1.TeamServiceClient)
	response, err := teamClient.GetTeamMembers(c.Request.Context(), &teamv1.GetTeamMembersRequest{
		TeamId: teamID,
	})

	if err != nil {
		h.logger.Error("Failed to get team members", "team_id", teamID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Team not found")
				return
			case codes.PermissionDenied:
				utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
				return
			case codes.Internal:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				return
			default:
				utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get team members")
		return
	}

	userClient := h.grpcClients.GetUserClient().(userv1.UserServiceClient)

	userIDs := make([]string, 0, len(response.Members))
	for _, member := range response.Members {
		userIDs = append(userIDs, member.UserId)
	}

	usersResponse, err := userClient.GetUsersByIds(c.Request.Context(), &userv1.GetUsersByIdsRequest{
		Ids: userIDs,
	})

	usersMap := make(map[string]models.User)
	if err == nil && usersResponse != nil {
		for _, user := range usersResponse.Users {
			usersMap[user.Id] = models.User{
				ID:        user.Id,
				Email:     user.Email,
				Username:  user.Username,
				CreatedAt: user.CreatedAt.AsTime(),
				UpdatedAt: user.UpdatedAt.AsTime(),
			}
		}
	} else {
		h.logger.Warn("Failed to get users info for team members", "team_id", teamID, "error", err)
	}

	members := make([]models.TeamMember, 0, len(response.Members))
	for _, member := range response.Members {
		teamMember := models.TeamMember{
			TeamID: member.TeamId,
			UserID: member.UserId,
			Role:   member.Role,
		}

		if user, exists := usersMap[member.UserId]; exists {
			teamMember.User = user
		} else {
			teamMember.User = models.User{
				ID: member.UserId,
			}
		}

		members = append(members, teamMember)
	}

	apiResponse := &models.GetTeamMembersResponse{
		Members: members,
	}

	utils.SuccessResponse(c, apiResponse)
}
