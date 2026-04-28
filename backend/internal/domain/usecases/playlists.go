package usecases

import (
	"net/http"

	"github.com/densmart/smart-stream/internal/adapters/db"
	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/models"
	"github.com/densmart/smart-stream/internal/domain/repo"
)

func CreatePlaylist(oltp repo.OltpRepo, data dto.CreatePlaylistDTO) (*models.Playlist, *UCError) {
	// create DB record
	playlist, dbErr := oltp.CreatePlaylist(data)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorDuplicateUniqueIndex {
			return nil, &UCError{
				Code:     211,
				HttpCode: http.StatusConflict,
				Message:  "Playlist already exists",
			}
		}
		return nil, &UCError{
			Code:     211,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return &playlist, nil
}

func RetrievePlaylist(oltp repo.OltpRepo, id string) (*models.Playlist, *UCError) {
	// retrieve DB record
	playlist, dbErr := oltp.RetrievePlaylist(id)
	if dbErr != nil {
		return nil, &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
		}
	}
	return &playlist, nil
}

func UpdatePlaylist(oltp repo.OltpRepo, id string, data dto.UpdatePlaylistDTO) (*models.Playlist, *UCError) {
	// Получаем текущую запись для проверки старого постера
	currentPlaylist, dbErr := oltp.RetrievePlaylist(id)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorNotFound {
			return nil, &UCError{
				Code:     215,
				HttpCode: http.StatusNotFound,
			}
		}
		return nil, &UCError{
			Code:     215,
			HttpCode: http.StatusInternalServerError,
		}
	}

	// Если обновляется постер И старый постер существует И они разные
	if data.Poster != nil && currentPlaylist.Poster != nil && *data.Poster != *currentPlaylist.Poster {
		// Удаляем старый постер (игнорируем ошибки, это не критично)
		if err := DeletePoster(*currentPlaylist.Poster); err != nil {
			// Логируем предупреждение, но продолжаем операцию
		}
	}

	// update DB record
	updatedPlaylist, dbErr := oltp.UpdatePlaylist(id, data)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorNotFound {
			return nil, &UCError{
				Code:     215,
				HttpCode: http.StatusNotFound,
			}
		}
		if dbErr.Code == db.DBErrorDuplicateUniqueIndex {
			return nil, &UCError{
				Code:     213,
				HttpCode: http.StatusBadRequest,
				Message:  "Playlist already exists",
			}
		}
		return nil, &UCError{
			Code:     213,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return &updatedPlaylist, nil
}

func DeletePlaylist(oltp repo.OltpRepo, id string) *UCError {
	// Получаем запись перед удалением чтобы удалить постер
	playlist, dbErr := oltp.RetrievePlaylist(id)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorNotFound {
			return &UCError{
				Code:     214,
				HttpCode: http.StatusNotFound,
			}
		}
		return &UCError{
			Code:     214,
			HttpCode: http.StatusInternalServerError,
		}
	}

	// Удаляем постер если есть (игнорируем ошибки)
	if playlist.Poster != nil {
		if err := DeletePoster(*playlist.Poster); err != nil {
			// Логируем предупреждение, но продолжаем операцию
		}
	}

	// Удаляем запись из БД
	dbErr = oltp.DeletePlaylist(id)
	if dbErr != nil {
		return &UCError{
			Code:     214,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func SearchPlaylists(oltp repo.OltpRepo, filter dto.SearchPlaylistsDTO) ([]models.Playlist, uint64, *UCError) {
	playlists, total, dbErr := oltp.SearchPlaylists(filter)
	if dbErr != nil {
		return nil, 0, &UCError{
			Code:     216,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return playlists, total, nil
}