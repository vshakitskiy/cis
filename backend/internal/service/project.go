package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/repository"
)

var ErrNotFound = errors.New("not found")

type ProjectService struct {
	projectRepo *repository.ProjectRepo
}

func NewProjectService(projectRepo *repository.ProjectRepo) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) Create(ctx context.Context, name string, description *string, ownerID int64, deadline *model.Date) (*model.Project, error) {
	project := &model.Project{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Deadline:    deadline,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return project, err
}

func (s *ProjectService) List(ctx context.Context) ([]*model.Project, error) {
	return s.projectRepo.List(ctx)
}

func (s *ProjectService) Update(ctx context.Context, id int64, name string, description *string, deadline *model.Date) (*model.Project, error) {
	project, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	project.Name = name
	project.Description = description
	project.Deadline = deadline

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Delete(ctx context.Context, id int64) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.projectRepo.Delete(ctx, id)
}
