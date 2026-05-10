package usecases

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/densmart/smart-stream/internal/adapters/db"
	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/models"
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/densmart/smart-stream/internal/domain/utils"
)

func CreateMedia(oltp repo.OltpRepo, data dto.CreateMediaDTO) (*models.Media, *UCError) {
	// Если Duration или Size не переданы, получаем их из файла
	if data.Duration == nil || data.Size == nil {
		// Построить полный путь к медиа-файлу
		mediaDir := viper.GetString("storage.media-dir")
		fullPath := filepath.Join(mediaDir, data.Path)

		// Получить метаданные из файла
		metadata, err := utils.GetMediaMetadata(fullPath)
		if err != nil {
			// Определяем тип ошибки
			if strings.Contains(err.Error(), "not found") {
				return nil, &UCError{
					Code:     230,
					HttpCode: http.StatusBadRequest,
					Message:  err.Error(),
				}
			}
			if strings.Contains(err.Error(), "unsupported") {
				return nil, &UCError{
					Code:     232,
					HttpCode: http.StatusBadRequest,
					Message:  err.Error(),
				}
			}
			return nil, &UCError{
				Code:     231,
				HttpCode: http.StatusInternalServerError,
				Message:  err.Error(),
			}
		}

		// Заполняем недостающие поля
		if data.Duration == nil {
			data.Duration = &metadata.Duration
		}
		if data.Size == nil {
			data.Size = &metadata.Size
		}
	}

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
	// Получаем текущую запись для проверки старого постера и пути
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

	// Если изменился Path И (Duration или Size не переданы), получаем метаданные из нового файла
	if data.Path != nil && *data.Path != currentMedia.Path && (data.Duration == nil || data.Size == nil) {
		mediaDir := viper.GetString("storage.media-dir")
		fullPath := filepath.Join(mediaDir, *data.Path)

		// Получить метаданные из файла
		metadata, err := utils.GetMediaMetadata(fullPath)
		if err != nil {
			// Определяем тип ошибки
			if strings.Contains(err.Error(), "not found") {
				return nil, &UCError{
					Code:     230,
					HttpCode: http.StatusBadRequest,
					Message:  err.Error(),
				}
			}
			if strings.Contains(err.Error(), "unsupported") {
				return nil, &UCError{
					Code:     232,
					HttpCode: http.StatusBadRequest,
					Message:  err.Error(),
				}
			}
			return nil, &UCError{
				Code:     231,
				HttpCode: http.StatusInternalServerError,
				Message:  err.Error(),
			}
		}

		// Заполняем недостающие поля из метаданных
		if data.Duration == nil {
			data.Duration = &metadata.Duration
		}
		if data.Size == nil {
			data.Size = &metadata.Size
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

// BrowseMediaFiles возвращает список файлов и папок в директории media
func BrowseMediaFiles(relativePath string) (*dto.BrowseMediaFilesResponse, *UCError) {
	// Получаем базовый путь из конфигурации
	mediaDir := viper.GetString("storage.media-dir")

	// Вызываем utility функцию для чтения директории
	items, parentPath, err := utils.BrowseDirectory(mediaDir, relativePath)
	if err != nil {
		// Определяем тип ошибки
		if strings.Contains(err.Error(), "not found") {
			return nil, &UCError{
				Code:     250,
				HttpCode: http.StatusNotFound,
				Message:  err.Error(),
			}
		}
		if strings.Contains(err.Error(), "invalid path") || strings.Contains(err.Error(), "not allowed") {
			return nil, &UCError{
				Code:     251,
				HttpCode: http.StatusBadRequest,
				Message:  err.Error(),
			}
		}
		return nil, &UCError{
			Code:     252,
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to browse directory",
		}
	}

	// Преобразуем результаты в DTO
	itemDTOs := make([]dto.FileSystemItemDTO, 0, len(items))
	for _, item := range items {
		itemDTO := dto.FileSystemItemDTO{
			Name: item.Name,
			Type: dto.FileSystemItemType(item.Type),
			Path: item.Path,
		}

		// Добавляем дополнительные поля только для файлов
		if item.Type == utils.FileSystemItemTypeFile {
			if item.Size != nil {
				itemDTO.Size = item.Size
			}
			if item.Format != nil {
				itemDTO.Format = item.Format
			}
			if item.ModifiedAt != nil {
				modifiedAtStr := item.ModifiedAt.Format("2006-01-02T15:04:05Z07:00")
				itemDTO.ModifiedAt = &modifiedAtStr
			}
		}

		itemDTOs = append(itemDTOs, itemDTO)
	}

	// Формируем response
	response := &dto.BrowseMediaFilesResponse{
		CurrentPath: relativePath,
		Items:       itemDTOs,
	}

	// ParentPath может быть nil для корневой директории
	if parentPath != "" {
		response.ParentPath = &parentPath
	}

	return response, nil
}