package models

type Media struct {
	BaseModel
	PlaylistID *string `db:"playlist_id" json:"playlist_id"`
	Name       string  `db:"name" json:"name"`
	Poster     *string `db:"poster" json:"poster"`
	Format     string  `db:"format" json:"format"`
	Path       string  `db:"path" json:"path"`
	Duration   int     `db:"duration" json:"duration"`
	Size       int     `db:"size" json:"size"`
	Order      int     `db:"order" json:"order"`
}

type Playlist struct {
	BaseModel
	Name   string  `db:"name" json:"name"`
	Type   string  `db:"type" json:"type"`
	Poster *string `db:"poster" json:"poster"`
}