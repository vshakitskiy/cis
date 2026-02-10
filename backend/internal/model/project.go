package model

import "time"

type Project struct {
	ID          int64     `json:"id"          db:"id"`
	Name        string    `json:"name"        db:"name"`
	Description *string   `json:"description" db:"description"`
	OwnerID     int64     `json:"owner_id"    db:"owner_id"`
	Deadline    *Date     `json:"deadline"    db:"deadline"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"  db:"updated_at"`
}
