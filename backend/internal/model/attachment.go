package model

import "time"

type Attachment struct {
	ID         int64     `json:"id"          db:"id"`
	TaskID     int64     `json:"task_id"     db:"task_id"`
	UploadedBy int64     `json:"uploaded_by" db:"uploaded_by"`
	Filename   string    `json:"filename"    db:"filename"`
	Filepath   string    `json:"-"           db:"filepath"`
	Size       int64     `json:"size"        db:"size"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
}
