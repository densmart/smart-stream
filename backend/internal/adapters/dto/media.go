package dto

type CreateMediaDTO struct {
	PlaylistID *string `json:"playlist_id"`
	Name       string  `json:"name" binding:"required"`
	Poster     *string `json:"poster"`
	Format     string  `json:"format" binding:"required"`
	Path       string  `json:"path" binding:"required"`
	Duration   int     `json:"duration" binding:"required,min=0"`
	Size       int     `json:"size" binding:"required,min=0"`
	Order      *int    `json:"order"`
}

type UpdateMediaDTO struct {
	PlaylistID *string `json:"playlist_id"`
	Name       *string `json:"name"`
	Poster     *string `json:"poster"`
	Format     *string `json:"format"`
	Path       *string `json:"path"`
	Duration   *int    `json:"duration"`
	Size       *int    `json:"size"`
	Order      *int    `json:"order"`
}

type SearchMediaDTO struct {
	BaseSearchRequestDTO
	ID         *string `form:"id"`
	PlaylistID *string `form:"playlist_id"`
	Name       *string `form:"name"`
	Format     *string `form:"format"`
}

type MediaDTO struct {
	ID         string  `json:"id"`
	CreatedAt  string  `json:"created_at"`
	PlaylistID *string `json:"playlist_id"`
	Name       string  `json:"name"`
	Poster     *string `json:"poster"`
	Format     string  `json:"format"`
	Path       string  `json:"path"`
	Duration   int     `json:"duration"`
	Size       int     `json:"size"`
	Order      int     `json:"order"`
}

type CreatePlaylistDTO struct {
	Name   string  `json:"name" binding:"required"`
	Type   string  `json:"type" binding:"required"`
	Poster *string `json:"poster"`
}

type UpdatePlaylistDTO struct {
	Name   *string `json:"name"`
	Type   *string `json:"type"`
	Poster *string `json:"poster"`
}

type SearchPlaylistsDTO struct {
	BaseSearchRequestDTO
	ID   *string `form:"id"`
	Name *string `form:"name"`
	Type *string `form:"type"`
}

type PlaylistDTO struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Poster    *string `json:"poster"`
}