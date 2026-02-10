package model

import "time"

type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	OwnerID     int64     `json:"owner_id"`
	Deadline    *string   `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
