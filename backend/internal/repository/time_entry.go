package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vshakitskiy/cis/internal/model"
)

type TimeEntryRepo struct {
	pool *pgxpool.Pool
}

func NewTimeEntryRepo(pool *pgxpool.Pool) *TimeEntryRepo {
	return &TimeEntryRepo{pool: pool}
}

func (r *TimeEntryRepo) Create(ctx context.Context, entry *model.TimeEntry) error {
	query := `
		INSERT INTO time_entries (task_id, user_id, minutes, description, date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		entry.TaskID, entry.UserID, entry.Minutes, entry.Description, entry.Date,
	).Scan(&entry.ID, &entry.CreatedAt)
}

func (r *TimeEntryRepo) GetByID(ctx context.Context, id int64) (*model.TimeEntry, error) {
	query := `
		SELECT id, task_id, user_id, minutes, description, date, created_at
		FROM time_entries WHERE id = $1
	`
	rows, _ := r.pool.Query(ctx, query, id)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.TimeEntry])
}

func (r *TimeEntryRepo) ListByTask(ctx context.Context, taskID int64) ([]*model.TimeEntry, error) {
	query := `
		SELECT id, task_id, user_id, minutes, description, date, created_at
		FROM time_entries WHERE task_id = $1
		ORDER BY date DESC, created_at DESC
	`
	rows, _ := r.pool.Query(ctx, query, taskID)
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[model.TimeEntry])
}

func (r *TimeEntryRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM time_entries WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *TimeEntryRepo) SumByProject(ctx context.Context, projectID int64) (int, error) {
	query := `
		SELECT COALESCE(SUM(te.minutes), 0)
		FROM time_entries te
		JOIN tasks t ON te.task_id = t.id
		WHERE t.project_id = $1
	`
	var total int
	err := r.pool.QueryRow(ctx, query, projectID).Scan(&total)
	return total, err
}
