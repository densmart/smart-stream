package dto

import "time"

type CreateClientDTO struct {
	Login    string  `json:"login"`
	Password string  `json:"password"`
	IsActive bool    `json:"is_active"`
	Email    *string `json:"email"`
}

type UpdateClientDTO struct {
	Login          *string `json:"login"`
	IsActive       *bool   `json:"is_active"`
	Email          *string `json:"email"`
	LastLoginAt    *time.Time
	CurrentVersion *string
	Password       *string
}

type SearchClientsDTO struct {
	BaseSearchRequestDTO
	ID              *string    `form:"id"`
	Login           *string    `form:"login"`
	Email           *string    `form:"email"`
	IsActive        *bool      `form:"is_active"`
	LastLoginAtFrom *time.Time `form:"last_login_at_from"`
	LastLoginAtTo   *time.Time `form:"last_login_at_to"`
	CurrentVersion  *string    `form:"current_version"`
}

type ClientDTO struct {
	ID             string  `json:"id"`
	CreatedAt      string  `json:"created_at"`
	Login          string  `json:"login"`
	Email          *string `json:"email"`
	IsActive       bool    `json:"is_active"`
	CurrentVersion *string `json:"current_version"`
	LastLoginAt    *string `json:"last_login_at"`
}

type ResetClientPasswordDTO struct {
	Password string `json:"password"`
}
