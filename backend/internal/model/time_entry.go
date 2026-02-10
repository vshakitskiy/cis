package model

import "time"

type TimeEntry struct {
	ID          int64     `json:"id"          db:"id"`
	TaskID      int64     `json:"task_id"     db:"task_id"`
	UserID      int64     `json:"user_id"     db:"user_id"`
	Minutes     int       `json:"minutes"     db:"minutes"`
	Description *string   `json:"description" db:"description"`
	Date        Date      `json:"date"        db:"date"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
}
