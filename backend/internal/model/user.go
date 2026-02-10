package model

import "time"

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleEmployee UserRole = "employee"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
