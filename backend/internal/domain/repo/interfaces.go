package repo

import (
	"github.com/densmart/smart-stream/internal/adapters/db"
	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/models"
)

type Clients interface {
	CreateClient(data dto.CreateClientDTO) (models.Client, *db.DBError)
	RetrieveClient(id string) (models.Client, *db.DBError)
	UpdateClient(id string, data dto.UpdateClientDTO) (models.Client, *db.DBError)
	DeleteClient(id string) *db.DBError
	SearchClients(data dto.SearchClientsDTO) ([]models.Client, uint64, *db.DBError)
}

type Media interface {
	CreateMedia(data dto.CreateMediaDTO) (models.Media, *db.DBError)
	RetrieveMedia(id string) (models.Media, *db.DBError)
	UpdateMedia(id string, data dto.UpdateMediaDTO) (models.Media, *db.DBError)
	DeleteMedia(id string) *db.DBError
	SearchMedia(data dto.SearchMediaDTO) ([]models.Media, uint64, *db.DBError)
}

type Playlists interface {
	CreatePlaylist(data dto.CreatePlaylistDTO) (models.Playlist, *db.DBError)
	RetrievePlaylist(id string) (models.Playlist, *db.DBError)
	UpdatePlaylist(id string, data dto.UpdatePlaylistDTO) (models.Playlist, *db.DBError)
	DeletePlaylist(id string) *db.DBError
	SearchPlaylists(data dto.SearchPlaylistsDTO) ([]models.Playlist, uint64, *db.DBError)
}
