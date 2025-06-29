package handlers

import (
	"net/http"
	"time"

	"api-gateway/internal/clients"
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	grpcClients *clients.GRPCClients
	logger      logger.Logger
}

func NewBoardHandler(grpcClients *clients.GRPCClients, log logger.Logger) *BoardHandler {
	return &BoardHandler{
		grpcClients: grpcClients,
		logger:      log,
	}
}

// CreateBoard godoc
// @Summary Create a new board
// @Description Create a new board for a team
// @Tags boards
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.CreateBoardRequest true "Create board request"
// @Success 201 {object} models.Response{data=models.Board}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/boards [post]
func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var req models.CreateBoardRequest
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

	h.logger.Info("Creating board", "name", req.Name, "team_id", req.TeamID, "created_by", req.CreatedBy)

	board := &models.Board{
		Name:        req.Name,
		Description: req.Description,
		TeamID:      req.TeamID,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	h.logger.Info("Board created successfully", "board_id", board.ID)
	utils.CreatedResponse(c, board)
}

// GetBoard godoc
// @Summary Get board with lists and tasks
// @Description Get board information including all lists and tasks
// @Tags boards
// @Produce json
// @Security Bearer
// @Param id path string true "Board ID"
// @Success 200 {object} models.Response{data=models.BoardWithLists}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/boards/{id} [get]
func (h *BoardHandler) GetBoard(c *gin.Context) {
	boardID := utils.GetParamID(c, "id")
	if boardID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Board ID is required")
		return
	}

	h.logger.Info("Getting board", "board_id", boardID)

	// TODO: Вызов gRPC сервиса досок
	boardWithLists := &models.BoardWithLists{
		Board: models.Board{
			ID: boardID,
		},
		Lists: []models.ListWithTasks{
			{
				List: models.List{
					ID:       "list_1",
					BoardID:  boardID,
					Name:     "To Do",
					Position: 1,
				},
				Tasks: []models.Task{
					{
						ID:          "task_1",
						ListID:      "list_1",
						Title:       "Example Task",
						Description: "Task description",
						Position:    1,
						CreatedBy:   "user_123",
					},
				},
			},
		},
	}

	utils.SuccessResponse(c, boardWithLists)
}

// UpdateBoard godoc
// @Summary Update board
// @Description Update board information
// @Tags boards
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Board ID"
// @Param request body models.UpdateBoardRequest true "Update board request"
// @Success 200 {object} models.Response{data=models.Board}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/boards/{id} [put]
func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	boardID := utils.GetParamID(c, "id")
	if boardID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Board ID is required")
		return
	}

	var req models.UpdateBoardRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Updating board", "board_id", boardID)

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	board := &models.Board{
		ID:          boardID,
		Name:        req.Name,
		Description: req.Description,
		UpdatedAt:   time.Now(),
	}

	h.logger.Info("Board updated successfully", "board_id", boardID)
	utils.SuccessResponse(c, board)
}

// DeleteBoard godoc
// @Summary Delete board
// @Description Delete board and all its contents
// @Tags boards
// @Produce json
// @Security Bearer
// @Param id path string true "Board ID"
// @Success 204 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/boards/{id} [delete]
func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	boardID := utils.GetParamID(c, "id")
	if boardID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Board ID is required")
		return
	}

	h.logger.Info("Deleting board", "board_id", boardID)

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	h.logger.Info("Board deleted successfully", "board_id", boardID)
	utils.NoContentResponse(c)
}

// GetUserBoards godoc
// @Summary Get user boards
// @Description Get all boards accessible to a user
// @Tags boards
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} models.Response{data=models.GetUserBoardsResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/users/{id}/boards [get]
func (h *BoardHandler) GetUserBoards(c *gin.Context) {
	userID := utils.GetParamID(c, "id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	h.logger.Info("Getting user boards", "user_id", userID)

	// TODO: Вызов gRPC сервиса досок
	boards := []models.Board{
		{
			ID:          "board_1",
			Name:        "Personal Board",
			Description: "My personal tasks",
			TeamID:      "team_1",
			CreatedBy:   userID,
		},
		{
			ID:          "board_2",
			Name:        "Team Board",
			Description: "Team collaboration",
			TeamID:      "team_2",
			CreatedBy:   "user_456",
		},
	}

	response := &models.GetUserBoardsResponse{
		Boards: boards,
	}

	utils.SuccessResponse(c, response)
}

// GetTeamBoards godoc
// @Summary Get team boards
// @Description Get all boards for a team
// @Tags boards
// @Produce json
// @Security Bearer
// @Param id path string true "Team ID"
// @Success 200 {object} models.Response{data=models.GetTeamBoardsResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/teams/{id}/boards [get]
func (h *BoardHandler) GetTeamBoards(c *gin.Context) {
	teamID := utils.GetParamID(c, "id")
	if teamID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Team ID is required")
		return
	}

	h.logger.Info("Getting team boards", "team_id", teamID)

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	boards := []models.Board{
		{
			ID:          "board_1",
			Name:        "Project Alpha",
			Description: "Alpha project board",
			TeamID:      teamID,
			CreatedBy:   "user_123",
		},
	}

	response := &models.GetTeamBoardsResponse{
		Boards: boards,
	}

	utils.SuccessResponse(c, response)
}

// CreateList godoc
// @Summary Create a list in board
// @Description Create a new list in a board
// @Tags boards
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Board ID"
// @Param request body models.CreateListRequest true "Create list request"
// @Success 201 {object} models.Response{data=models.List}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/boards/{id}/lists [post]
func (h *BoardHandler) CreateList(c *gin.Context) {
	boardID := utils.GetParamID(c, "id")
	if boardID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Board ID is required")
		return
	}

	var req models.CreateListRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	req.BoardID = boardID

	h.logger.Info("Creating list", "board_id", boardID, "name", req.Name)

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	list := &models.List{
		ID:        "list_123",
		BoardID:   boardID,
		Name:      req.Name,
		Position:  req.Position,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	h.logger.Info("List created successfully", "list_id", list.ID)
	utils.CreatedResponse(c, list)
}

// UpdateList godoc
// @Summary Update list
// @Description Update list information
// @Tags boards
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "List ID"
// @Param request body models.UpdateListRequest true "Update list request"
// @Success 200 {object} models.Response{data=models.List}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/lists/{id} [put]
func (h *BoardHandler) UpdateList(c *gin.Context) {
	listID := utils.GetParamID(c, "id")
	if listID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "List ID is required")
		return
	}

	var req models.UpdateListRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Updating list", "list_id", listID)

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	list := &models.List{
		ID:        listID,
		Name:      req.Name,
		UpdatedAt: time.Now(),
	}

	h.logger.Info("List updated successfully", "list_id", listID)
	utils.SuccessResponse(c, list)
}

// DeleteList godoc
// @Summary Delete list
// @Description Delete list and all its tasks
// @Tags boards
// @Produce json
// @Security Bearer
// @Param id path string true "List ID"
// @Success 204 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/lists/{id} [delete]
func (h *BoardHandler) DeleteList(c *gin.Context) {
	listID := utils.GetParamID(c, "id")
	if listID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "List ID is required")
		return
	}

	h.logger.Info("Deleting list", "list_id", listID)

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	h.logger.Info("List deleted successfully", "list_id", listID)
	utils.NoContentResponse(c)
}

// ReorderLists godoc
// @Summary Reorder lists in board
// @Description Change the order of lists in a board
// @Tags boards
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Board ID"
// @Param request body models.ReorderListsRequest true "Reorder lists request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/boards/{id}/lists/reorder [put]
func (h *BoardHandler) ReorderLists(c *gin.Context) {
	boardID := utils.GetParamID(c, "id")
	if boardID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Board ID is required")
		return
	}

	var req models.ReorderListsRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Reordering lists", "board_id", boardID, "positions_count", len(req.Positions))

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	h.logger.Info("Lists reordered successfully", "board_id", boardID)
	utils.SuccessResponse(c, gin.H{"message": "Lists reordered successfully"})
}

// GetBoardLabels godoc
// @Summary Get board labels
// @Description Get all labels for a board
// @Tags boards
// @Produce json
// @Security Bearer
// @Param id path string true "Board ID"
// @Success 200 {object} models.Response{data=models.GetBoardLabelsResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/boards/{id}/labels [get]
func (h *BoardHandler) GetBoardLabels(c *gin.Context) {
	boardID := utils.GetParamID(c, "id")
	if boardID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Board ID is required")
		return
	}

	h.logger.Info("Getting board labels", "board_id", boardID)

	// TODO: Проверка прав доступа и вызов gRPC сервиса досок
	labels := []models.Label{
		{
			ID:      "label_1",
			BoardID: boardID,
			Name:    "Bug",
			Color:   "#ff0000",
		},
		{
			ID:      "label_2",
			BoardID: boardID,
			Name:    "Feature",
			Color:   "#00ff00",
		},
	}

	response := &models.GetBoardLabelsResponse{
		Labels: labels,
	}

	utils.SuccessResponse(c, response)
}
