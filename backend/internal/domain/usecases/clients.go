package usecases

import (
	"net/http"
	"time"

	"github.com/densmart/smart-stream/internal/adapters/db"
	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/models"
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/densmart/smart-stream/internal/domain/utils"
)

func CreateClient(oltp repo.OltpRepo, data dto.CreateClientDTO) (*models.Client, *UCError) {
	// create password hash
	pwdHash, err := utils.GeneratePasswordHash(data.Password)
	if err != nil {
		return nil, &UCError{
			Code:     157,
			HttpCode: http.StatusInternalServerError,
		}
	}
	data.Password = pwdHash

	// create DB record
	client, dbErr := oltp.CreateClient(data)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorDuplicateUniqueIndex {
			return nil, &UCError{
				Code:     151,
				HttpCode: http.StatusConflict,
				Message:  "Client already exists",
			}
		}
		return nil, &UCError{
			Code:     151,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return &client, nil
}

func RetrieveClient(oltp repo.OltpRepo, id string) (*models.Client, *UCError) {
	// retrieve DB record
	client, dbErr := oltp.RetrieveClient(id)
	if dbErr != nil {
		return nil, &UCError{
			Code:     152,
			HttpCode: http.StatusNotFound,
		}
	}
	return &client, nil
}

func UpdateClient(oltp repo.OltpRepo, id string, data dto.UpdateClientDTO) (*models.Client, *UCError) {
	// update DB record
	updatedClient, dbErr := oltp.UpdateClient(id, data)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorNotFound {
			return nil, &UCError{
				Code:     155,
				HttpCode: http.StatusNotFound,
			}
		}
		if dbErr.Code == db.DBErrorDuplicateUniqueIndex {
			return nil, &UCError{
				Code:     153,
				HttpCode: http.StatusBadRequest,
				Message:  "Client already exists",
			}
		}
		return nil, &UCError{
			Code:     153,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return &updatedClient, nil
}

func DeleteClient(oltp repo.OltpRepo, id string) *UCError {
	dbErr := oltp.DeleteClient(id)
	if dbErr != nil {
		return &UCError{
			Code:     154,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func SearchClients(oltp repo.OltpRepo, filter dto.SearchClientsDTO) ([]models.Client, uint64, *UCError) {
	clients, total, dbErr := oltp.SearchClients(filter)
	if dbErr != nil {
		return nil, 0, &UCError{
			Code:     156,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return clients, total, nil
}

func ResetClientPassword(oltp repo.OltpRepo, id string, data dto.ResetClientPasswordDTO) (*models.Client, *UCError) {
	// create password hash
	pwdHash, err := utils.GeneratePasswordHash(data.Password)
	if err != nil {
		return nil, &UCError{
			Code:     157,
			HttpCode: http.StatusInternalServerError,
		}
	}

	updateData := dto.UpdateClientDTO{
		Password:    &pwdHash,
		LastLoginAt: new(time.Time),
	}

	updatedClient, dbErr := oltp.UpdateClient(id, updateData)
	if dbErr != nil {
		return nil, &UCError{
			Code:     153,
			HttpCode: http.StatusInternalServerError,
		}
	}

	return &updatedClient, nil
}
