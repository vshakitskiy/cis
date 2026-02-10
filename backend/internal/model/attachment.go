package model

import "time"

type Attachment struct {
	ID         int64     `json:"id"`
	TaskID     int64     `json:"task_id"`
	UploadedBy int64     `json:"uploaded_by"`
	Filename   string    `json:"filename"`
	Filepath   string    `json:"-"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
}
