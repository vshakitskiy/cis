package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vshakitskiy/cis/internal/model"
)

type ProjectRepo struct {
	pool *pgxpool.Pool
}

func NewProjectRepo(pool *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{pool: pool}
}

func (r *ProjectRepo) Create(ctx context.Context, project *model.Project) error {
	query := `
		INSERT INTO projects (name, description, owner_id, deadline)
		VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, updated_at
  `

	return r.pool.QueryRow(ctx, query,
		project.Name, project.Description, project.OwnerID, project.Deadline,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *ProjectRepo) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	query := `
		SELECT id, name, description, owner_id, deadline, created_at, updated_at
		FROM projects WHERE id = $1
  `

	rows, _ := r.pool.Query(ctx, query, id)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Project])
}

func (r *ProjectRepo) List(ctx context.Context) ([]*model.Project, error) {
	query := `
		SELECT id, name, description, owner_id, deadline, created_at, updated_at
		FROM projects ORDER BY created_at DESC
  `

	rows, _ := r.pool.Query(ctx, query)
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Project])
}

func (r *ProjectRepo) Update(ctx context.Context, project *model.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2, deadline = $3, updated_at = now()
		WHERE id = $4
		RETURNING updated_at
  `

	return r.pool.QueryRow(ctx, query,
		project.Name, project.Description, project.Deadline, project.ID,
	).Scan(&project.UpdatedAt)
}

func (r *ProjectRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
