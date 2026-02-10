package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/repository"
)

type AttachmentService struct {
	attachmentRepo *repository.AttachmentRepo
	taskRepo       *repository.TaskRepo
	uploadDir      string
}

func NewAttachmentService(attachmentRepo *repository.AttachmentRepo, taskRepo *repository.TaskRepo, uploadDir string) *AttachmentService {
	os.MkdirAll(uploadDir, 0o755)
	return &AttachmentService{
		attachmentRepo: attachmentRepo,
		taskRepo:       taskRepo,
		uploadDir:      uploadDir,
	}
}

func (s *AttachmentService) Create(ctx context.Context, taskID, userID int64, filename string, size int64, file io.Reader) (*model.Attachment, error) {
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	storedName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	storedPath := filepath.Join(s.uploadDir, storedName)

	dst, err := os.Create(storedPath)
	if err != nil {
		return nil, fmt.Errorf("creating file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(storedPath)
		return nil, fmt.Errorf("writing file: %w", err)
	}

	attachment := &model.Attachment{
		TaskID:     taskID,
		UploadedBy: userID,
		Filename:   filename,
		Filepath:   storedPath,
		Size:       size,
	}

	if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
		os.Remove(storedPath)
		return nil, err
	}

	return attachment, nil
}

func (s *AttachmentService) ListByTask(ctx context.Context, taskID int64) ([]*model.Attachment, error) {
	return s.attachmentRepo.ListByTask(ctx, taskID)
}

func (s *AttachmentService) GetByID(ctx context.Context, id int64) (*model.Attachment, error) {
	attachment, err := s.attachmentRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return attachment, err
}

func (s *AttachmentService) Delete(ctx context.Context, id int64) error {
	attachment, err := s.attachmentRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if err := s.attachmentRepo.Delete(ctx, id); err != nil {
		return err
	}

	os.Remove(attachment.Filepath)
	return nil
}
