package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vshakitskiy/cis/internal/model"
)

type AttachmentRepo struct {
	pool *pgxpool.Pool
}

func NewAttachmentRepo(pool *pgxpool.Pool) *AttachmentRepo {
	return &AttachmentRepo{pool: pool}
}

func (r *AttachmentRepo) Create(ctx context.Context, attachment *model.Attachment) error {
	query := `
		INSERT INTO attachments (task_id, uploaded_by, filename, filepath, size)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		attachment.TaskID, attachment.UploadedBy, attachment.Filename, attachment.Filepath, attachment.Size,
	).Scan(&attachment.ID, &attachment.CreatedAt)
}

func (r *AttachmentRepo) GetByID(ctx context.Context, id int64) (*model.Attachment, error) {
	query := `
		SELECT id, task_id, uploaded_by, filename, filepath, size, created_at
		FROM attachments WHERE id = $1
	`
	rows, _ := r.pool.Query(ctx, query, id)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Attachment])
}

func (r *AttachmentRepo) ListByTask(ctx context.Context, taskID int64) ([]*model.Attachment, error) {
	query := `
		SELECT id, task_id, uploaded_by, filename, filepath, size, created_at
		FROM attachments WHERE task_id = $1
		ORDER BY created_at DESC
	`
	rows, _ := r.pool.Query(ctx, query, taskID)
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.Attachment])
}

func (r *AttachmentRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM attachments WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
