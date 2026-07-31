package dto

type CreateMediaDTO struct {
	PlaylistID *string `json:"playlist_id"`
	Name       string  `json:"name" binding:"required"`
	Poster     *string `json:"poster"`
	Format     string  `json:"format" binding:"required"`
	Path       string  `json:"path" binding:"required"`
	Duration   *int    `json:"duration" binding:"omitempty,min=0"`
	Size       *uint64 `json:"size" binding:"omitempty,min=0"`
	Order      *int    `json:"order"`
}

type UpdateMediaDTO struct {
	PlaylistID *string `json:"playlist_id"`
	Name       *string `json:"name"`
	Poster     *string `json:"poster"`
	Format     *string `json:"format"`
	Path       *string `json:"path"`
	Duration   *int    `json:"duration"`
	Size       *uint64 `json:"size"`
	Order      *int    `json:"order"`
}

type SearchMediaDTO struct {
	BaseSearchRequestDTO
	ID             *string `form:"id"`
	PlaylistID     *string `form:"playlist_id"`
	Name           *string `form:"name"`
	Format         *string `form:"format"`
	OnlyUnassigned *bool   `form:"only_unassigned"` // Только медиа без плейлистов
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
	Size       uint64  `json:"size"`
	Order      int     `json:"order"`
}

type CreatePlaylistDTO struct {
	Name     string  `json:"name" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	Poster   *string `json:"poster"`
	ParentID *string `json:"parent_id"`
}

type UpdatePlaylistDTO struct {
	Name     *string `json:"name"`
	Type     *string `json:"type"`
	Poster   *string `json:"poster"`
	ParentID *string `json:"parent_id"`
}

type SearchPlaylistsDTO struct {
	BaseSearchRequestDTO
	ID                 *string `form:"id"`
	Name               *string `form:"name"`
	Type               *string `form:"type"`
	ParentID           *string `form:"parent_id"`
	IncludeAllHierarchy *bool   `form:"include_all_hierarchy"` // If true, load all playlists regardless of parent_id
}

type PlaylistDTO struct {
	ID          string  `json:"id"`
	CreatedAt   string  `json:"created_at"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Poster      *string `json:"poster"`
	ParentID    *string `json:"parent_id"`
	HasChildren bool    `json:"has_children"`
}

// Playlist search autocomplete
type PlaylistSearchDTO struct {
	Name *string `form:"name"`
}

type PlaylistSearchResultDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// File browser DTOs
type FileSystemItemType string

const (
	FileSystemItemTypeFile      FileSystemItemType = "file"
	FileSystemItemTypeDirectory FileSystemItemType = "directory"
)

type FileSystemItemDTO struct {
	Name       string             `json:"name"`
	Type       FileSystemItemType `json:"type"`
	Path       string             `json:"path"`
	Size       *uint64            `json:"size,omitempty"`        // Only for files
	Format     *string            `json:"format,omitempty"`      // Only for files
	ModifiedAt *string            `json:"modified_at,omitempty"` // Only for files
}

type BrowseMediaFilesResponse struct {
	CurrentPath string              `json:"current_path"`
	ParentPath  *string             `json:"parent_path"`
	Items       []FileSystemItemDTO `json:"items"`
}

// Playlist media DTOs
type AddMediaToPlaylistDTO struct {
	MediaID string `json:"media_id" binding:"required"`
	Order   int    `json:"order" binding:"min=0"`
}

type UpdateMediaOrderDTO struct {
	Order int `json:"order" binding:"min=0"`
}

// Batch Playlist media DTOs
type BatchAddMediaToPlaylistDTO struct {
	Media []struct {
		MediaID string `json:"media_id" binding:"required"`
		Order   int    `json:"order" binding:"min=0"`
	} `json:"media" binding:"required,min=1,dive"`
}

type BatchRemoveMediaFromPlaylistDTO struct {
	MediaIDs []string `json:"media_ids" binding:"required,min=1"`
}

type BatchUpdateMediaOrderDTO struct {
	Updates []struct {
		MediaID string `json:"media_id" binding:"required"`
		Order   int    `json:"order" binding:"min=0"`
	} `json:"updates" binding:"required,min=1,dive"`
}