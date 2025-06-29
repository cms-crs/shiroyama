package handlers

import (
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"net/http"

	"api-gateway/internal/clients"
	"api-gateway/internal/models"
	"api-gateway/internal/utils"
	"api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskHandler struct {
	grpcClients *clients.GRPCClients
	logger      logger.Logger
}

func NewTaskHandler(grpcClients *clients.GRPCClients, log logger.Logger) *TaskHandler {
	return &TaskHandler{
		grpcClients: grpcClients,
		logger:      log,
	}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task in a list
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.CreateTaskRequest true "Create task request"
// @Success 201 {object} models.Response{data=models.Task}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
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

	h.logger.Info("Creating task", "title", req.Title, "list_id", req.ListID, "created_by", req.CreatedBy)

	taskClient := h.grpcClients.GetTaskClient().(taskv1.TaskServiceClient)

	grpcReq := &taskv1.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
	}

	if req.ListID != "" {
		grpcReq.ListId = req.ListID
	}
	if req.CreatedBy != "" {
		grpcReq.CreatedBy = req.CreatedBy
	}
	if req.Position > 0 {
		grpcReq.Position = int32(req.Position)
	}
	if req.DueDate.IsZero() {
		grpcReq.DueDate = timestamppb.New(req.DueDate)
	}

	response, err := taskClient.CreateTask(c.Request.Context(), grpcReq)

	if err != nil {
		h.logger.Error("Failed to create task", "error", err)

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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	task := &models.Task{
		ID:          response.Id,
		Title:       response.Title,
		Description: response.Description,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
	}

	if response.ListId != "" {
		task.ListID = response.ListId
	}
	if response.CreatedBy != "" {
		task.CreatedBy = response.CreatedBy
	}
	if response.Position > 0 {
		task.Position = response.Position
	}
	if response.DueDate != nil {
		dueDate := response.DueDate.AsTime()
		task.DueDate = dueDate
	}

	h.logger.Info("Task created successfully", "task_id", task.ID)
	utils.CreatedResponse(c, task)
}

// GetTask godoc
// @Summary Get task by ID
// @Description Get task information by ID
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Success 200 {object} models.Response{data=models.Task}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := utils.GetParamID(c, "id")
	if taskID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Task ID is required")
		return
	}

	h.logger.Info("Getting task", "task_id", taskID)

	taskClient := h.grpcClients.GetTaskClient().(taskv1.TaskServiceClient)
	response, err := taskClient.GetTask(c.Request.Context(), &taskv1.GetTaskRequest{
		Id: taskID,
	})

	if err != nil {
		h.logger.Error("Failed to get task", "task_id", taskID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get task")
		return
	}

	task := &models.Task{
		ID:          response.Id,
		Title:       response.Title,
		Description: response.Description,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
	}

	if response.ListId != "" {
		task.ListID = response.ListId
	}
	if response.CreatedBy != "" {
		task.CreatedBy = response.CreatedBy
	}
	if response.Position > 0 {
		task.Position = response.Position
	}

	if response.DueDate != nil {
		dueDate := response.DueDate.AsTime()
		task.DueDate = dueDate
	}

	utils.SuccessResponse(c, task)
}

// UpdateTask godoc
// @Summary Update task
// @Description Update task information
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Param request body models.UpdateTaskRequest true "Update task request"
// @Success 200 {object} models.Response{data=models.Task}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := utils.GetParamID(c, "id")
	if taskID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Task ID is required")
		return
	}

	var req models.UpdateTaskRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	h.logger.Info("Updating task", "task_id", taskID)

	taskClient := h.grpcClients.GetTaskClient().(taskv1.TaskServiceClient)

	grpcReq := &taskv1.UpdateTaskRequest{
		Id:          taskID,
		Title:       req.Title,
		Description: req.Description,
	}

	if req.DueDate.IsZero() {
		grpcReq.DueDate = timestamppb.New(req.DueDate)
	}

	response, err := taskClient.UpdateTask(c.Request.Context(), grpcReq)

	if err != nil {
		h.logger.Error("Failed to update task", "task_id", taskID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update task")
		return
	}

	task := &models.Task{
		ID:          response.Id,
		Title:       response.Title,
		Description: response.Description,
		CreatedAt:   response.CreatedAt.AsTime(),
		UpdatedAt:   response.UpdatedAt.AsTime(),
	}

	if response.ListId != "" {
		task.ListID = response.ListId
	}

	if response.DueDate != nil {
		dueDate := response.DueDate.AsTime()
		task.DueDate = dueDate
	}

	h.logger.Info("Task updated successfully", "task_id", taskID)
	utils.SuccessResponse(c, task)
}

// DeleteTask godoc
// @Summary Delete task
// @Description Delete task
// @Tags tasks
// @Produce json
// @Security Bearer
// @Param id path string true "Task ID"
// @Success 204 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := utils.GetParamID(c, "id")
	if taskID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Task ID is required")
		return
	}

	h.logger.Info("Deleting task", "task_id", taskID)

	taskClient := h.grpcClients.GetTaskClient().(taskv1.TaskServiceClient)
	_, err := taskClient.DeleteTask(c.Request.Context(), &taskv1.DeleteTaskRequest{
		Id: taskID,
	})

	if err != nil {
		h.logger.Error("Failed to delete task", "task_id", taskID, "error", err)

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.ErrorResponse(c, http.StatusBadRequest, st.Message())
				return
			case codes.NotFound:
				utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	h.logger.Info("Task deleted successfully", "task_id", taskID)
	utils.NoContentResponse(c)
}
