package model

import "time"

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleEmployee UserRole = "employee"
)

type User struct {
	ID        int64     `json:"id"         db:"id"`
	Email     string    `json:"email"      db:"email"`
	Password  string    `json:"-"          db:"password"`
	Name      string    `json:"name"       db:"name"`
	Role      UserRole  `json:"role"       db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
