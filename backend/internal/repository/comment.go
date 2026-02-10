package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vshakitskiy/cis/internal/model"
)

type CommentRepo struct {
	pool *pgxpool.Pool
}

func NewCommentRepo(pool *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{pool: pool}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (task_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		comment.TaskID, comment.UserID, comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
}

func (r *CommentRepo) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	query := `
		SELECT id, task_id, user_id, content, created_at, updated_at
		FROM comments WHERE id = $1
	`
	rows, _ := r.pool.Query(ctx, query, id)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Comment])
}

func (r *CommentRepo) ListByTask(ctx context.Context, taskID int64) ([]*model.Comment, error) {
	query := `
		SELECT id, task_id, user_id, content, created_at, updated_at
		FROM comments WHERE task_id = $1
		ORDER BY created_at ASC
	`
	rows, _ := r.pool.Query(ctx, query, taskID)
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Comment])
}

func (r *CommentRepo) Update(ctx context.Context, comment *model.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, updated_at = now()
		WHERE id = $2
		RETURNING updated_at
	`
	return r.pool.QueryRow(ctx, query,
		comment.Content, comment.ID,
	).Scan(&comment.UpdatedAt)
}

func (r *CommentRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
