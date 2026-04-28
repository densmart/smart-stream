package usecases

import (
	"net/http"

	"github.com/densmart/smart-stream/internal/adapters/db"
	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/models"
	"github.com/densmart/smart-stream/internal/domain/repo"
)

func CreateMedia(oltp repo.OltpRepo, data dto.CreateMediaDTO) (*models.Media, *UCError) {
	// create DB record
	media, dbErr := oltp.CreateMedia(data)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorDuplicateUniqueIndex {
			return nil, &UCError{
				Code:     201,
				HttpCode: http.StatusConflict,
				Message:  "Media already exists",
			}
		}
		return nil, &UCError{
			Code:     201,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return &media, nil
}

func RetrieveMedia(oltp repo.OltpRepo, id string) (*models.Media, *UCError) {
	// retrieve DB record
	media, dbErr := oltp.RetrieveMedia(id)
	if dbErr != nil {
		return nil, &UCError{
			Code:     202,
			HttpCode: http.StatusNotFound,
		}
	}
	return &media, nil
}

func UpdateMedia(oltp repo.OltpRepo, id string, data dto.UpdateMediaDTO) (*models.Media, *UCError) {
	// Получаем текущую запись для проверки старого постера
	currentMedia, dbErr := oltp.RetrieveMedia(id)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorNotFound {
			return nil, &UCError{
				Code:     205,
				HttpCode: http.StatusNotFound,
			}
		}
		return nil, &UCError{
			Code:     205,
			HttpCode: http.StatusInternalServerError,
		}
	}

	// Если обновляется постер И старый постер существует И они разные
	if data.Poster != nil && currentMedia.Poster != nil && *data.Poster != *currentMedia.Poster {
		// Удаляем старый постер (игнорируем ошибки, это не критично)
		if err := DeletePoster(*currentMedia.Poster); err != nil {
			// Логируем предупреждение, но продолжаем операцию
		}
	}

	// update DB record
	updatedMedia, dbErr := oltp.UpdateMedia(id, data)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorNotFound {
			return nil, &UCError{
				Code:     205,
				HttpCode: http.StatusNotFound,
			}
		}
		if dbErr.Code == db.DBErrorDuplicateUniqueIndex {
			return nil, &UCError{
				Code:     203,
				HttpCode: http.StatusBadRequest,
				Message:  "Media already exists",
			}
		}
		return nil, &UCError{
			Code:     203,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return &updatedMedia, nil
}

func DeleteMedia(oltp repo.OltpRepo, id string) *UCError {
	// Получаем запись перед удалением чтобы удалить постер
	media, dbErr := oltp.RetrieveMedia(id)
	if dbErr != nil {
		if dbErr.Code == db.DBErrorNotFound {
			return &UCError{
				Code:     204,
				HttpCode: http.StatusNotFound,
			}
		}
		return &UCError{
			Code:     204,
			HttpCode: http.StatusInternalServerError,
		}
	}

	// Удаляем постер если есть (игнорируем ошибки)
	if media.Poster != nil {
		if err := DeletePoster(*media.Poster); err != nil {
			// Логируем предупреждение, но продолжаем операцию
		}
	}

	// Удаляем запись из БД
	dbErr = oltp.DeleteMedia(id)
	if dbErr != nil {
		return &UCError{
			Code:     204,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func SearchMedia(oltp repo.OltpRepo, filter dto.SearchMediaDTO) ([]models.Media, uint64, *UCError) {
	media, total, dbErr := oltp.SearchMedia(filter)
	if dbErr != nil {
		return nil, 0, &UCError{
			Code:     206,
			HttpCode: http.StatusInternalServerError,
		}
	}
	return media, total, nil
}