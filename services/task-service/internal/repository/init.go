package repository

import (
	"context"
	"fmt"
	"time"

	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// GORM Models
type Task struct {
	ID          string     `gorm:"primaryKey;column:id"`
	ListID      string     `gorm:"column:list_id;not null"`
	Title       string     `gorm:"column:title;not null"`
	Description string     `gorm:"column:description"`
	Position    int32      `gorm:"column:position;not null"`
	DueDate     *time.Time `gorm:"column:due_date"`
	CreatedBy   string     `gorm:"column:created_by;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"`

	// Associations
	AssignedUsers []TaskAssignment `gorm:"foreignKey:TaskID"`
	Labels        []Label          `gorm:"many2many:task_labels;"`
}

func (Task) TableName() string {
	return "tasks"
}

type TaskAssignment struct {
	ID         uint       `gorm:"primaryKey"`
	TaskID     string     `gorm:"column:task_id;not null"`
	UserID     string     `gorm:"column:user_id;not null"`
	AssignedAt time.Time  `gorm:"column:assigned_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;index"`
}

func (TaskAssignment) TableName() string {
	return "task_assignments"
}

type Label struct {
	ID        string     `gorm:"primaryKey;column:id"`
	BoardID   string     `gorm:"column:board_id;not null"`
	Name      string     `gorm:"column:name;not null"`
	Color     string     `gorm:"column:color;not null"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (Label) TableName() string {
	return "labels"
}

type TaskLabel struct {
	TaskID    string     `gorm:"column:task_id;primaryKey"`
	LabelID   string     `gorm:"column:label_id;primaryKey"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (TaskLabel) TableName() string {
	return "task_labels"
}

// Repository
type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *taskv1.Task) (*taskv1.Task, error) {
	dbTask := &Task{
		ListID:      task.ListId,
		Title:       task.Title,
		Description: task.Description,
		Position:    task.Position,
		CreatedBy:   task.CreatedBy,
	}

	if task.DueDate != nil {
		dueDate := task.DueDate.AsTime()
		dbTask.DueDate = &dueDate
	}

	if err := r.db.WithContext(ctx).Create(dbTask).Error; err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return r.dbTaskToProto(dbTask), nil
}

func (r *TaskRepository) GetTask(ctx context.Context, taskID string) (*taskv1.Task, error) {
	var dbTask Task

	err := r.db.WithContext(ctx).
		Preload("AssignedUsers", "deleted_at IS NULL").
		Preload("Labels", "deleted_at IS NULL").
		Where("id = ? AND deleted_at IS NULL", taskID).
		First(&dbTask).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return r.dbTaskToProto(&dbTask), nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.Task, error) {
	updates := map[string]interface{}{
		"title":       req.Title,
		"description": req.Description,
		"updated_at":  time.Now(),
	}

	if req.DueDate != nil {
		updates["due_date"] = req.DueDate.AsTime()
	} else {
		updates["due_date"] = nil
	}

	err := r.db.WithContext(ctx).
		Model(&Task{}).
		Where("id = ? AND deleted_at IS NULL", req.Id).
		Updates(updates).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Check if task was actually updated
	var count int64
	r.db.WithContext(ctx).Model(&Task{}).Where("id = ? AND deleted_at IS NULL", req.Id).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("task not found")
	}

	return r.GetTask(ctx, req.Id)
}

func (r *TaskRepository) DeleteTask(ctx context.Context, taskID string) error {
	result := r.db.WithContext(ctx).
		Model(&Task{}).
		Where("id = ? AND deleted_at IS NULL", taskID).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to delete task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) MoveTask(ctx context.Context, req *taskv1.MoveTaskRequest) (*taskv1.Task, error) {
	updates := map[string]interface{}{
		"list_id":    req.ToListId,
		"position":   req.Position,
		"updated_at": time.Now(),
	}

	err := r.db.WithContext(ctx).
		Model(&Task{}).
		Where("id = ? AND deleted_at IS NULL", req.TaskId).
		Updates(updates).Error

	if err != nil {
		return nil, fmt.Errorf("failed to move task: %w", err)
	}

	// Check if task was actually updated
	var count int64
	r.db.WithContext(ctx).Model(&Task{}).Where("id = ? AND deleted_at IS NULL", req.TaskId).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("task not found")
	}

	return r.GetTask(ctx, req.TaskId)
}

func (r *TaskRepository) AssignUser(ctx context.Context, taskID, userID string) (*taskv1.Task, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Check if user is already assigned
		var count int64
		tx.WithContext(ctx).
			Model(&TaskAssignment{}).
			Where("task_id = ? AND user_id = ? AND deleted_at IS NULL", taskID, userID).
			Count(&count)

		if count == 0 {
			// Create assignment
			assignment := &TaskAssignment{
				TaskID:     taskID,
				UserID:     userID,
				AssignedAt: time.Now(),
			}

			if err := tx.WithContext(ctx).Create(assignment).Error; err != nil {
				return fmt.Errorf("failed to assign user: %w", err)
			}
		}

		// Update task updated_at
		if err := tx.WithContext(ctx).
			Model(&Task{}).
			Where("id = ? AND deleted_at IS NULL", taskID).
			Update("updated_at", time.Now()).Error; err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return updated task
	return r.GetTask(ctx, taskID)
}

func (r *TaskRepository) UnassignUser(ctx context.Context, taskID, userID string) (*taskv1.Task, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Remove assignment (soft delete)
		if err := tx.WithContext(ctx).
			Model(&TaskAssignment{}).
			Where("task_id = ? AND user_id = ? AND deleted_at IS NULL", taskID, userID).
			Update("deleted_at", time.Now()).Error; err != nil {
			return fmt.Errorf("failed to unassign user: %w", err)
		}

		// Update task updated_at
		if err := tx.WithContext(ctx).
			Model(&Task{}).
			Where("id = ? AND deleted_at IS NULL", taskID).
			Update("updated_at", time.Now()).Error; err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return updated task
	return r.GetTask(ctx, taskID)
}

func (r *TaskRepository) GetTasksForLists(ctx context.Context, listIDs []string) ([]*taskv1.Task, error) {
	if len(listIDs) == 0 {
		return []*taskv1.Task{}, nil
	}

	var dbTasks []Task
	err := r.db.WithContext(ctx).
		Preload("AssignedUsers", "deleted_at IS NULL").
		Preload("Labels", "deleted_at IS NULL").
		Where("list_id IN ? AND deleted_at IS NULL", listIDs).
		Order("list_id, position").
		Find(&dbTasks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for lists: %w", err)
	}

	tasks := make([]*taskv1.Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		task := r.dbTaskToProto(&dbTask)
		tasks[i] = task
	}

	return tasks, nil
}

func (r *TaskRepository) GetTasksForUser(ctx context.Context, userID string) ([]*taskv1.Task, error) {
	var dbTasks []Task

	err := r.db.WithContext(ctx).
		Distinct().
		Preload("AssignedUsers", "deleted_at IS NULL").
		Preload("Labels", "deleted_at IS NULL").
		Joins("INNER JOIN task_assignments ta ON tasks.id = ta.task_id").
		Where("ta.user_id = ? AND tasks.deleted_at IS NULL AND ta.deleted_at IS NULL", userID).
		Order("tasks.created_at DESC").
		Find(&dbTasks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for user: %w", err)
	}

	tasks := make([]*taskv1.Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		task := r.dbTaskToProto(&dbTask)
		tasks[i] = task
	}

	return tasks, nil
}

func (r *TaskRepository) dbTaskToProto(dbTask *Task) *taskv1.Task {
	task := &taskv1.Task{
		Id:          dbTask.ID,
		ListId:      dbTask.ListID,
		Title:       dbTask.Title,
		Description: dbTask.Description,
		Position:    dbTask.Position,
		CreatedBy:   dbTask.CreatedBy,
		CreatedAt:   timestamppb.New(dbTask.CreatedAt),
		UpdatedAt:   timestamppb.New(dbTask.UpdatedAt),
	}

	if dbTask.DueDate != nil {
		task.DueDate = timestamppb.New(*dbTask.DueDate)
	}

	return task
}

func (r *TaskRepository) getTaskAssignedUsers(ctx context.Context, taskID string) ([]string, error) {
	var assignments []TaskAssignment
	err := r.db.WithContext(ctx).
		Where("task_id = ? AND deleted_at IS NULL", taskID).
		Order("assigned_at").
		Find(&assignments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get assigned users: %w", err)
	}

	userIDs := make([]string, len(assignments))
	for i, assignment := range assignments {
		userIDs[i] = assignment.UserID
	}

	return userIDs, nil
}
