package models

import "time"

type Client struct {
	BaseModel
	Login          string     `db:"login" json:"login"`
	Password       string     `db:"password" json:"password"`
	IsActive       bool       `db:"is_active" json:"is_active"`
	Email          *string    `db:"email" json:"email"`
	CurrentVersion *string    `db:"current_version" json:"current_version"`
	LastLoginAt    *time.Time `db:"last_login_at" json:"last_login_at"`
}
