package model

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type TaskPriority string

const (
	PriorityLow      TaskPriority = "low"
	PriorityMedium   TaskPriority = "medium"
	PriorityHigh     TaskPriority = "high"
	PriorityCritical TaskPriority = "critical"
)

type Task struct {
	ID          int64        `json:"id"          db:"id"`
	ProjectID   int64        `json:"project_id"  db:"project_id"`
	AssigneeID  *int64       `json:"assignee_id" db:"assignee_id"`
	Title       string       `json:"title"       db:"title"`
	Description *string      `json:"description" db:"description"`
	Status      TaskStatus   `json:"status"      db:"status"`
	Priority    TaskPriority `json:"priority"    db:"priority"`
	Deadline    *Date        `json:"deadline"    db:"deadline"`
	CreatedAt   time.Time    `json:"created_at"  db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"  db:"updated_at"`
}
