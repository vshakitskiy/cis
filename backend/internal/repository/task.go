package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vshakitskiy/cis/internal/model"
)

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{pool: pool}
}

func (r *TaskRepo) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (project_id, assignee_id, title, description, status, priority, deadline)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, created_at, updated_at
  `

	return r.pool.QueryRow(ctx, query,
		task.ProjectID, task.AssigneeID, task.Title, task.Description,
		task.Status, task.Priority, task.Deadline,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepo) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	query := `
  	SELECT id, project_id, assignee_id, title, description, status, priority, deadline, created_at, updated_at
   	FROM tasks WHERE id = $1
  `

	rows, _ := r.pool.Query(ctx, query, id)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Task])
}

func (r *TaskRepo) ListByProject(ctx context.Context, projectID int64) ([]*model.Task, error) {
	query := `
    SELECT id, project_id, assignee_id, title, description, status, priority, deadline, created_at, updated_at
    FROM tasks WHERE project_id = $1
    ORDER BY created_at DESC
  `

	rows, _ := r.pool.Query(ctx, query, projectID)
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Task])
}

func (r *TaskRepo) Update(ctx context.Context, task *model.Task) error {
	query := `
    UPDATE tasks
    SET assignee_id = $1, title = $2, description = $3, status = $4, priority = $5, deadline = $6, updated_at = now()
    WHERE id = $7
    RETURNING updated_at
  `

	return r.pool.QueryRow(ctx, query,
		task.AssigneeID, task.Title, task.Description, task.Status, task.Priority, task.Deadline, task.ID,
	).Scan(&task.UpdatedAt)
}

func (r *TaskRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
