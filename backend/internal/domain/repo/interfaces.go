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
