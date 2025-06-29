package handlers

import (
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"net/http"

	"api-gateway/internal/clients"
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	response, err := boardClient.CreateBoard(c.Request.Context(), &boardv1.CreateBoardRequest{
		Name:        req.Name,
		Description: req.Description,
		TeamId:      req.TeamID,
		CreatedBy:   req.CreatedBy,
	})

	if err != nil {
		h.logger.Error("Failed to create board", "error", err)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create board")
		return
	}

	board := &models.Board{
		ID:          response.Id,
		Name:        response.Name,
		Description: response.Description,
		TeamID:      response.TeamId,
		CreatedBy:   response.CreatedBy,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
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
	//boardID := utils.GetParamID(c, "id")
	//if boardID == "" {
	//	utils.ErrorResponse(c, http.StatusBadRequest, "Board ID is required")
	//	return
	//}
	//
	//h.logger.Info("Getting board", "board_id", boardID)
	//
	//boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	//response, err := boardClient.GetBoard(c.Request.Context(), &boardv1.GetBoardRequest{
	//	Id: boardID,
	//})
	//
	//if err != nil {
	//	h.logger.Error("Failed to get board", "board_id", boardID, "error", err)
	//
	//	if st, ok := status.FromError(err); ok {
	//		switch st.Code() {
	//		case codes.InvalidArgument:
	//			utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
	//			return
	//		case codes.NotFound:
	//			utils.ErrorResponse(c, http.StatusNotFound, "Board not found")
	//			return
	//		case codes.PermissionDenied:
	//			utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
	//			return
	//		case codes.Internal:
	//			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	//			return
	//		default:
	//			utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
	//			return
	//		}
	//	}
	//	utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get board")
	//	return
	//}
	//
	//boardWithLists := &models.BoardWithLists{
	//	Board: models.Board{
	//		ID:          response.Board.Id,
	//		Name:        response.Board.Name,
	//		Description: response.Board.Description,
	//		TeamID:      response.Board.TeamId,
	//		CreatedBy:   response.Board.CreatedBy,
	//		CreatedAt:   response.Board.CreatedAt.AsTime(),
	//		UpdatedAt:   response.Board.UpdatedAt.AsTime(),
	//	},
	//	Lists: make([]models.ListWithTasks, len(response.Lists)),
	//}
	//
	//for i, list := range response.Lists {
	//	listWithTasks := models.ListWithTasks{
	//		List: models.List{
	//			ID:        list.List.Id,
	//			BoardID:   list.List.BoardId,
	//			Name:      list.List.Name,
	//			Position:  int(list.List.Position),
	//			CreatedAt: list.List.CreatedAt.AsTime(),
	//			UpdatedAt: list.List.UpdatedAt.AsTime(),
	//		},
	//		Tasks: make([]models.Task, len(list.Tasks)),
	//	}
	//
	//	for j, task := range list.Tasks {
	//		listWithTasks.Tasks[j] = models.Task{
	//			ID:          task.Id,
	//			ListID:      task.ListId,
	//			Title:       task.Title,
	//			Description: task.Description,
	//			Position:    task.Position,
	//			CreatedBy:   task.CreatedBy,
	//			CreatedAt:   task.CreatedAt.AsTime(),
	//			UpdatedAt:   task.UpdatedAt.AsTime(),
	//		}
	//		if task.DueDate != nil {
	//			dueDate := task.DueDate.AsTime()
	//			listWithTasks.Tasks[j].DueDate = &dueDate
	//		}
	//	}
	//
	//	boardWithLists.Lists[i] = listWithTasks
	//}
	//
	//utils.SuccessResponse(c, boardWithLists)
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

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	response, err := boardClient.UpdateBoard(c.Request.Context(), &boardv1.UpdateBoardRequest{
		Id:          boardID,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		h.logger.Error("Failed to update board", "board_id", boardID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Board not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update board")
		return
	}

	// Конвертация protobuf модели в модель API
	board := &models.Board{
		ID:          response.Id,
		Name:        response.Name,
		Description: response.Description,
		TeamID:      response.TeamId,
		CreatedBy:   response.CreatedBy,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
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

	// TODO: Добавить метод DeleteBoard в gRPC сервис
	// boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	// _, err := boardClient.DeleteBoard(c.Request.Context(), &boardv1.DeleteBoardRequest{
	//     Id: boardID,
	// })

	// if err != nil {
	//     h.logger.Error("Failed to delete board", "board_id", boardID, "error", err)
	//
	//     if st, ok := status.FromError(err); ok {
	//         switch st.Code() {
	//         case codes.InvalidArgument:
	//             utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
	//             return
	//         case codes.NotFound:
	//             utils.ErrorResponse(c, http.StatusNotFound, "Board not found")
	//             return
	//         case codes.PermissionDenied:
	//             utils.ErrorResponse(c, http.StatusForbidden, "Permission denied")
	//             return
	//         case codes.Internal:
	//             utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	//             return
	//         default:
	//             utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error occurred")
	//             return
	//         }
	//     }
	//     utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete board")
	//     return
	// }

	// DeleteBoard не реализован в gRPC сервисе
	utils.ErrorResponse(c, http.StatusNotImplemented, "Delete board not implemented yet")
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

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	response, err := boardClient.GetUserBoards(c.Request.Context(), &boardv1.GetUserBoardsRequest{
		UserId: userID,
	})

	if err != nil {
		h.logger.Error("Failed to get user boards", "user_id", userID, "error", err)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user boards")
		return
	}

	boards := make([]models.Board, len(response.Boards))
	for i, board := range response.Boards {
		boards[i] = models.Board{
			ID:          board.Id,
			Name:        board.Name,
			Description: board.Description,
			TeamID:      board.TeamId,
			CreatedBy:   board.CreatedBy,
			CreatedAt:   board.CreatedAt.AsTime(),
			UpdatedAt:   board.UpdatedAt.AsTime(),
		}
	}

	apiResponse := &models.GetUserBoardsResponse{
		Boards: boards,
	}

	utils.SuccessResponse(c, apiResponse)
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

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	response, err := boardClient.GetTeamBoards(c.Request.Context(), &boardv1.GetTeamBoardsRequest{
		TeamId: teamID,
	})

	if err != nil {
		h.logger.Error("Failed to get team boards", "team_id", teamID, "error", err)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get team boards")
		return
	}

	boards := make([]models.Board, len(response.Boards))
	for i, board := range response.Boards {
		boards[i] = models.Board{
			ID:          board.Id,
			Name:        board.Name,
			Description: board.Description,
			TeamID:      board.TeamId,
			CreatedBy:   board.CreatedBy,
			CreatedAt:   board.CreatedAt.AsTime(),
			UpdatedAt:   board.UpdatedAt.AsTime(),
		}
	}

	apiResponse := &models.GetTeamBoardsResponse{
		Boards: boards,
	}

	utils.SuccessResponse(c, apiResponse)
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

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	response, err := boardClient.CreateList(c.Request.Context(), &boardv1.CreateListRequest{
		BoardId:  boardID,
		Name:     req.Name,
		Position: int32(req.Position),
	})

	if err != nil {
		h.logger.Error("Failed to create list", "board_id", boardID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Board not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create list")
		return
	}

	list := &models.List{
		ID:        response.Id,
		BoardID:   response.BoardId,
		Name:      response.Name,
		Position:  int32(int(response.Position)),
		CreatedAt: response.CreatedAt.AsTime(),
		UpdatedAt: response.UpdatedAt.AsTime(),
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

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	response, err := boardClient.UpdateList(c.Request.Context(), &boardv1.UpdateListRequest{
		Id:   listID,
		Name: req.Name,
	})

	if err != nil {
		h.logger.Error("Failed to update list", "list_id", listID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "List not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update list")
		return
	}

	list := &models.List{
		ID:        response.Id,
		BoardID:   response.BoardId,
		Name:      response.Name,
		Position:  int32(int(response.Position)),
		CreatedAt: response.CreatedAt.AsTime(),
		UpdatedAt: response.UpdatedAt.AsTime(),
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

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	_, err := boardClient.DeleteList(c.Request.Context(), &boardv1.DeleteListRequest{
		Id: listID,
	})

	if err != nil {
		h.logger.Error("Failed to delete list", "list_id", listID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "List not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete list")
		return
	}

	utils.ErrorResponse(c, http.StatusNotImplemented, "Delete list not implemented yet")
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

	positions := make([]*boardv1.ListPosition, len(req.Positions))
	for i, pos := range req.Positions {
		positions[i] = &boardv1.ListPosition{
			ListId:   pos.ListID,
			Position: int32(pos.Position),
		}
	}

	boardClient := h.grpcClients.GetBoardClient().(boardv1.BoardServiceClient)
	_, err := boardClient.ReorderLists(c.Request.Context(), &boardv1.ReorderListsRequest{
		BoardId:   boardID,
		Positions: positions,
	})

	if err != nil {
		h.logger.Error("Failed to reorder lists", "board_id", boardID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Board not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reorder lists")
		return
	}

	h.logger.Info("Lists reordered successfully", "board_id", boardID)
	utils.SuccessResponse(c, gin.H{"message": "Lists reordered successfully"})
}
