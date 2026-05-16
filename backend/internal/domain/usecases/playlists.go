package usecases

import (
	"net/http"

	"github.com/densmart/smart-stream/internal/adapters/db"
	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/models"
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/densmart/smart-stream/internal/domain/utils"
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
	// Получаем запись перед удалением чтобы удалить постер и обновить родителя
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

	// Обновляем has_children у родителя если он был
	if playlist.ParentID != nil {
		if dbErr := oltp.UpdateParentHasChildren(*playlist.ParentID); dbErr != nil {
			// Логируем предупреждение, но не возвращаем ошибку
			// т.к. плейлист уже удален
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

// GetPlaylistMedia возвращает медиа плейлиста с сортировкой по order
func GetPlaylistMedia(oltp repo.OltpRepo, playlistID string, filter dto.SearchMediaDTO) ([]models.Media, uint64, *UCError) {
	// Проверяем существование плейлиста
	_, dbErr := oltp.RetrievePlaylist(playlistID)
	if dbErr != nil {
		return nil, 0, &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
			Message:  "Playlist not found",
		}
	}

	// Устанавливаем фильтр по playlist_id
	filter.PlaylistID = &playlistID

	// Устанавливаем сортировку по order если не указана другая
	if filter.Order == nil {
		orderBy := "order"
		filter.Order = &orderBy
	}

	// Получаем медиа
	media, total, dbErr := oltp.SearchMedia(filter)
	if dbErr != nil {
		return nil, 0, &UCError{
			Code:     260,
			HttpCode: http.StatusInternalServerError,
		}
	}

	return media, total, nil
}

// AddMediaToPlaylist добавляет медиа в плейлист
func AddMediaToPlaylist(oltp repo.OltpRepo, playlistID string, data dto.AddMediaToPlaylistDTO) *UCError {
	// Проверяем существование плейлиста
	playlist, dbErr := oltp.RetrievePlaylist(playlistID)
	if dbErr != nil {
		return &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
			Message:  "Playlist not found",
		}
	}

	// Проверяем что плейлист не имеет дочерних плейлистов
	if playlist.HasChildren {
		return &UCError{
			Code:     267,
			HttpCode: http.StatusBadRequest,
			Message:  "Cannot add media to playlist with children",
		}
	}

	// Проверяем существование медиа
	media, dbErr := oltp.RetrieveMedia(data.MediaID)
	if dbErr != nil {
		return &UCError{
			Code:     202,
			HttpCode: http.StatusNotFound,
			Message:  "Media not found",
		}
	}

	// Проверяем что медиа не в плейлисте
	if media.PlaylistID != nil {
		if *media.PlaylistID == playlistID {
			return &UCError{
				Code:     262,
				HttpCode: http.StatusBadRequest,
				Message:  "Media already in this playlist",
			}
		}
		return &UCError{
			Code:     263,
			HttpCode: http.StatusBadRequest,
			Message:  "Media already in another playlist",
		}
	}

	// Добавляем медиа в плейлист
	dbErr = oltp.AddMediaToPlaylist(data.MediaID, playlistID, data.Order)
	if dbErr != nil {
		return &UCError{
			Code:     261,
			HttpCode: http.StatusInternalServerError,
		}
	}

	return nil
}

// RemoveMediaFromPlaylist удаляет медиа из плейлиста
func RemoveMediaFromPlaylist(oltp repo.OltpRepo, playlistID string, mediaID string) *UCError {
	// Проверяем существование плейлиста
	_, dbErr := oltp.RetrievePlaylist(playlistID)
	if dbErr != nil {
		return &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
			Message:  "Playlist not found",
		}
	}

	// Проверяем существование медиа
	media, dbErr := oltp.RetrieveMedia(mediaID)
	if dbErr != nil {
		return &UCError{
			Code:     202,
			HttpCode: http.StatusNotFound,
			Message:  "Media not found",
		}
	}

	// Проверяем что медиа в плейлисте
	if media.PlaylistID == nil || *media.PlaylistID != playlistID {
		return &UCError{
			Code:     265,
			HttpCode: http.StatusBadRequest,
			Message:  "Media not in this playlist",
		}
	}

	// Удаляем медиа из плейлиста
	dbErr = oltp.RemoveMediaFromPlaylist(mediaID)
	if dbErr != nil {
		return &UCError{
			Code:     264,
			HttpCode: http.StatusInternalServerError,
		}
	}

	return nil
}

// UpdateMediaOrder обновляет order медиа в плейлисте
func UpdateMediaOrder(oltp repo.OltpRepo, playlistID string, mediaID string, order int) *UCError {
	// Проверяем существование плейлиста
	_, dbErr := oltp.RetrievePlaylist(playlistID)
	if dbErr != nil {
		return &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
			Message:  "Playlist not found",
		}
	}

	// Проверяем существование медиа
	media, dbErr := oltp.RetrieveMedia(mediaID)
	if dbErr != nil {
		return &UCError{
			Code:     202,
			HttpCode: http.StatusNotFound,
			Message:  "Media not found",
		}
	}

	// Проверяем что медиа в плейлисте
	if media.PlaylistID == nil || *media.PlaylistID != playlistID {
		return &UCError{
			Code:     265,
			HttpCode: http.StatusBadRequest,
			Message:  "Media not in this playlist",
		}
	}

	// Обновляем order
	dbErr = oltp.UpdateMediaOrder(mediaID, order)
	if dbErr != nil {
		return &UCError{
			Code:     266,
			HttpCode: http.StatusInternalServerError,
		}
	}

	return nil
}

// BatchAddMediaToPlaylist добавляет несколько медиа в плейлист за один запрос
func BatchAddMediaToPlaylist(oltp repo.OltpRepo, playlistID string, data dto.BatchAddMediaToPlaylistDTO) *UCError {
	// Проверяем существование плейлиста
	playlist, dbErr := oltp.RetrievePlaylist(playlistID)
	if dbErr != nil {
		return &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
			Message:  "Playlist not found",
		}
	}

	// Проверяем что плейлист не имеет дочерних плейлистов
	if playlist.HasChildren {
		return &UCError{
			Code:     267,
			HttpCode: http.StatusBadRequest,
			Message:  "Cannot add media to playlist with children",
		}
	}

	// Валидируем все медиа перед добавлением
	for _, item := range data.Media {
		media, dbErr := oltp.RetrieveMedia(item.MediaID)
		if dbErr != nil {
			return &UCError{
				Code:     202,
				HttpCode: http.StatusNotFound,
				Message:  "Media not found: " + item.MediaID,
			}
		}

		// Проверяем что медиа не в плейлисте
		if media.PlaylistID != nil {
			if *media.PlaylistID == playlistID {
				return &UCError{
					Code:     262,
					HttpCode: http.StatusBadRequest,
					Message:  "Media already in this playlist: " + item.MediaID,
				}
			}
			return &UCError{
				Code:     263,
				HttpCode: http.StatusBadRequest,
				Message:  "Media already in another playlist: " + item.MediaID,
			}
		}
	}

	// Добавляем все медиа
	for _, item := range data.Media {
		dbErr = oltp.AddMediaToPlaylist(item.MediaID, playlistID, item.Order)
		if dbErr != nil {
			return &UCError{
				Code:     261,
				HttpCode: http.StatusInternalServerError,
				Message:  "Failed to add media: " + item.MediaID,
			}
		}
	}

	return nil
}

// BatchRemoveMediaFromPlaylist удаляет несколько медиа из плейлиста за один запрос
func BatchRemoveMediaFromPlaylist(oltp repo.OltpRepo, playlistID string, data dto.BatchRemoveMediaFromPlaylistDTO) *UCError {
	// Проверяем существование плейлиста
	_, dbErr := oltp.RetrievePlaylist(playlistID)
	if dbErr != nil {
		return &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
			Message:  "Playlist not found",
		}
	}

	// Валидируем все медиа перед удалением
	for _, mediaID := range data.MediaIDs {
		media, dbErr := oltp.RetrieveMedia(mediaID)
		if dbErr != nil {
			return &UCError{
				Code:     202,
				HttpCode: http.StatusNotFound,
				Message:  "Media not found: " + mediaID,
			}
		}

		// Проверяем что медиа в плейлисте
		if media.PlaylistID == nil || *media.PlaylistID != playlistID {
			return &UCError{
				Code:     265,
				HttpCode: http.StatusBadRequest,
				Message:  "Media not in this playlist: " + mediaID,
			}
		}
	}

	// Удаляем все медиа
	for _, mediaID := range data.MediaIDs {
		dbErr = oltp.RemoveMediaFromPlaylist(mediaID)
		if dbErr != nil {
			return &UCError{
				Code:     264,
				HttpCode: http.StatusInternalServerError,
				Message:  "Failed to remove media: " + mediaID,
			}
		}
	}

	return nil
}

// BatchUpdateMediaOrder обновляет order нескольких медиа в плейлисте за один запрос
func BatchUpdateMediaOrder(oltp repo.OltpRepo, playlistID string, data dto.BatchUpdateMediaOrderDTO) *UCError {
	// Проверяем существование плейлиста
	_, dbErr := oltp.RetrievePlaylist(playlistID)
	if dbErr != nil {
		return &UCError{
			Code:     212,
			HttpCode: http.StatusNotFound,
			Message:  "Playlist not found",
		}
	}

	// Валидируем все медиа перед обновлением
	for _, update := range data.Updates {
		media, dbErr := oltp.RetrieveMedia(update.MediaID)
		if dbErr != nil {
			return &UCError{
				Code:     202,
				HttpCode: http.StatusNotFound,
				Message:  "Media not found: " + update.MediaID,
			}
		}

		// Проверяем что медиа в плейлисте
		if media.PlaylistID == nil || *media.PlaylistID != playlistID {
			return &UCError{
				Code:     265,
				HttpCode: http.StatusBadRequest,
				Message:  "Media not in this playlist: " + update.MediaID,
			}
		}
	}

	// Обновляем order для всех медиа
	for _, update := range data.Updates {
		dbErr = oltp.UpdateMediaOrder(update.MediaID, update.Order)
		if dbErr != nil {
			return &UCError{
				Code:     266,
				HttpCode: http.StatusInternalServerError,
				Message:  "Failed to update media order: " + update.MediaID,
			}
		}
	}

	return nil
}

// SearchPlaylistsAutocomplete возвращает список плейлистов для автокомплита
func SearchPlaylistsAutocomplete(oltp repo.OltpRepo, filter dto.PlaylistSearchDTO) ([]dto.PlaylistSearchResultDTO, *UCError) {
	searchFilter := dto.SearchPlaylistsDTO{
		BaseSearchRequestDTO: dto.BaseSearchRequestDTO{
			Limit: utils.Ptr(uint(50)),
		},
		Name:                filter.Name,
		IncludeAllHierarchy: utils.Ptr(true),
	}

	playlists, _, dbErr := oltp.SearchPlaylists(searchFilter)
	if dbErr != nil {
		return nil, &UCError{
			Code:     216,
			HttpCode: http.StatusInternalServerError,
		}
	}

	results := make([]dto.PlaylistSearchResultDTO, len(playlists))
	for i, playlist := range playlists {
		results[i] = dto.PlaylistSearchResultDTO{
			ID:   playlist.ID,
			Name: playlist.Name,
		}
	}

	return results, nil
}