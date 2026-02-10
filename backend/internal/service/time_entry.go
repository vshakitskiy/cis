package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/repository"
)

type TimeEntryService struct {
	timeEntryRepo *repository.TimeEntryRepo
	taskRepo      *repository.TaskRepo
}

func NewTimeEntryService(timeEntryRepo *repository.TimeEntryRepo, taskRepo *repository.TaskRepo) *TimeEntryService {
	return &TimeEntryService{
		timeEntryRepo: timeEntryRepo,
		taskRepo:      taskRepo,
	}
}

func (s *TimeEntryService) Create(ctx context.Context, taskID, userID int64, minutes int, description *string, date model.Date) (*model.TimeEntry, error) {
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	entry := &model.TimeEntry{
		TaskID:      taskID,
		UserID:      userID,
		Minutes:     minutes,
		Description: description,
		Date:        date,
	}

	if err := s.timeEntryRepo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *TimeEntryService) ListByTask(ctx context.Context, taskID int64) ([]*model.TimeEntry, error) {
	return s.timeEntryRepo.ListByTask(ctx, taskID)
}

func (s *TimeEntryService) Delete(ctx context.Context, id int64) error {
	_, err := s.timeEntryRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	return s.timeEntryRepo.Delete(ctx, id)
}
