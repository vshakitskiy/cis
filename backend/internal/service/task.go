package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/repository"
)

type TaskService struct {
	taskRepo    *repository.TaskRepo
	projectRepo *repository.ProjectRepo
}

func NewTaskService(taskRepo *repository.TaskRepo, projectRepo *repository.ProjectRepo) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskService) Create(ctx context.Context, projectID int64, assigneeID *int64, title string, description *string, priority model.TaskPriority,
	deadline *model.Date) (*model.Task, error) {
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	task := &model.Task{
		ProjectID:   projectID,
		AssigneeID:  assigneeID,
		Title:       title,
		Description: description,
		Status:      model.StatusTodo,
		Priority:    priority,
		Deadline:    deadline,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return task, err
}

func (s *TaskService) ListByProject(ctx context.Context, projectID int64) ([]*model.Task, error) {
	return s.taskRepo.ListByProject(ctx, projectID)
}

func (s *TaskService) Update(ctx context.Context, id int64, assigneeID *int64, title string, description *string, status model.TaskStatus, priority model.TaskPriority, deadline *model.Date) (*model.Task, error) {
	task, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.AssigneeID = assigneeID
	task.Title = title
	task.Description = description
	task.Status = status
	task.Priority = priority
	task.Deadline = deadline

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.taskRepo.Delete(ctx, id)
}
