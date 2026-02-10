package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepo struct {
	pool *pgxpool.Pool
}

func NewReportRepo(pool *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{pool: pool}
}

type ProjectReport struct {
	ProjectID      int64   `json:"project_id"`
	ProjectName    string  `json:"project_name"`
	TotalTasks     int     `json:"total_tasks"`
	CompletedTasks int     `json:"completed_tasks"`
	CompletionPct  float64 `json:"completion_pct"`
	TotalMinutes   int     `json:"total_minutes"`
}

func (r *ReportRepo) GetProjectReport(ctx context.Context, projectID int64) (*ProjectReport, error) {
	query := `
		SELECT
			p.id AS project_id,
			p.name AS project_name,
			COUNT(t.id) AS total_tasks,
			COUNT(t.id) FILTER (WHERE t.status = 'done') AS completed_tasks,
			CASE
				WHEN COUNT(t.id) = 0 THEN 0
				ELSE ROUND(COUNT(t.id) FILTER (WHERE t.status = 'done')::numeric / COUNT(t.id) * 100, 1)
			END AS completion_pct,
			COALESCE((
				SELECT SUM(te.minutes)
				FROM time_entries te
				JOIN tasks t2 ON te.task_id = t2.id
				WHERE t2.project_id = p.id
			), 0) AS total_minutes
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.id = $1
		GROUP BY p.id, p.name
	`

	report := &ProjectReport{}
	err := r.pool.QueryRow(ctx, query, projectID).Scan(
		&report.ProjectID,
		&report.ProjectName,
		&report.TotalTasks,
		&report.CompletedTasks,
		&report.CompletionPct,
		&report.TotalMinutes,
	)
	if err != nil {
		return nil, err
	}

	return report, nil
}
