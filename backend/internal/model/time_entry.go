package model

import "time"

type TimeEntry struct {
	ID          int64     `json:"id"`
	TaskID      int64     `json:"task_id"`
	UserID      int64     `json:"user_id"`
	Minutes     int       `json:"minutes"`
	Description *string   `json:"description"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}
