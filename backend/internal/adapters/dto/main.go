package dto

import "time"

type BaseSearchRequestDTO struct {
	RawURL            string
	CreatedAtFrom     time.Time `form:"created_at_from"`
	CreatedAtTo       time.Time `form:"created_at_to"`
	Order             *string   `form:"order"`
	Offset            *uint     `form:"offset"`
	Limit             *uint     `form:"limit"`
	AvailableEntities []string
}

type PaginationInfo struct {
	Total   uint   `json:"total"`
	Page    uint   `json:"page"`
	PerPage uint   `json:"per_page"`
	Pages   uint   `json:"pages"`
	Next    string `json:"next"`
	Prev    string `json:"prev"`
}

type AutocompleteDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
