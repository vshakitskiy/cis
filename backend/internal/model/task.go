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
	ID          int64        `json:"id"`
	ProjectID   int64        `json:"project_id"`
	AssigneeID  *int64       `json:"assignee_id"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	Deadline    *string      `json:"deadline"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
