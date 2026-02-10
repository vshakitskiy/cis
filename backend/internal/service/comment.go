package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepo
	taskRepo    *repository.TaskRepo
}

func NewCommentService(commentRepo *repository.CommentRepo, taskRepo *repository.TaskRepo) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		taskRepo:    taskRepo,
	}
}

func (s *CommentService) Create(ctx context.Context, taskID, userID int64, content string) (*model.Comment, error) {
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	comment := &model.Comment{
		TaskID:  taskID,
		UserID:  userID,
		Content: content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) ListByTask(ctx context.Context, taskID int64) ([]*model.Comment, error) {
	return s.commentRepo.ListByTask(ctx, taskID)
}

func (s *CommentService) Update(ctx context.Context, id int64, content string) (*model.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	comment.Content = content

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) Delete(ctx context.Context, id int64) error {
	_, err := s.commentRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	return s.commentRepo.Delete(ctx, id)
}
