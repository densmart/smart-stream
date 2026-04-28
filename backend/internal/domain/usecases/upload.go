package usecases

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/densmart/smart-stream/pkg/logger"
)

var allowedPosterExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

var allowedPosterMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
}

type UploadedFile struct {
	Path     string
	Filename string
}

// UploadPoster validates and saves a poster image file
func UploadPoster(file *multipart.FileHeader) (*UploadedFile, *UCError) {
	// Проверяем размер файла
	maxSizeMB := viper.GetInt64("storage.max-poster-size-mb")
	maxSize := maxSizeMB * 1024 * 1024 // Конвертируем MB в байты
	if file.Size > maxSize {
		return nil, &UCError{
			Code:     221,
			HttpCode: http.StatusBadRequest,
			Message:  fmt.Sprintf("File size exceeds maximum allowed size of %d MB", maxSizeMB),
		}
	}

	// Проверяем MIME-тип
	contentType := file.Header.Get("Content-Type")
	if !allowedPosterMimeTypes[contentType] {
		return nil, &UCError{
			Code:     222,
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid file type. Allowed types: jpg, jpeg, png, webp",
		}
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedPosterExtensions[ext] {
		return nil, &UCError{
			Code:     223,
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid file extension. Allowed extensions: .jpg, .jpeg, .png, .webp",
		}
	}

	// Генерируем уникальное имя файла
	filename := uuid.New().String() + ext

	// Получаем путь к директории для постеров из конфига
	postersDir := viper.GetString("storage.posters-dir")

	// Создаем директорию, если её нет
	if err := os.MkdirAll(postersDir, 0755); err != nil {
		logger.Errorf("[usecases] UploadPoster error creating directory: %s", err.Error())
		return nil, &UCError{
			Code:     224,
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to create storage directory",
		}
	}

	// Полный путь к файлу
	fullPath := filepath.Join(postersDir, filename)

	// Открываем загруженный файл
	src, err := file.Open()
	if err != nil {
		logger.Errorf("[usecases] UploadPoster error opening file: %s", err.Error())
		return nil, &UCError{
			Code:     225,
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to open uploaded file",
		}
	}
	defer src.Close()

	// Создаем файл на диске
	dst, err := os.Create(fullPath)
	if err != nil {
		logger.Errorf("[usecases] UploadPoster error creating file: %s", err.Error())
		return nil, &UCError{
			Code:     226,
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to save file",
		}
	}
	defer dst.Close()

	// Копируем содержимое
	if _, err = io.Copy(dst, src); err != nil {
		logger.Errorf("[usecases] UploadPoster error copying file: %s", err.Error())
		// Удаляем частично созданный файл
		os.Remove(fullPath)
		return nil, &UCError{
			Code:     227,
			HttpCode: http.StatusInternalServerError,
			Message:  "Failed to save file",
		}
	}

	logger.Infof("[usecases] UploadPoster successfully uploaded file: %s", filename)

	// Возвращаем относительный путь к файлу
	relativePath := "/static/posters/" + filename

	return &UploadedFile{
		Path:     relativePath,
		Filename: filename,
	}, nil
}

// DeletePoster removes a poster file from disk
func DeletePoster(filename string) error {
	if filename == "" {
		return nil
	}

	// Извлекаем только имя файла из пути (на случай если передан полный путь)
	baseName := filepath.Base(filename)
	// Удаляем префикс /static/posters/ если он есть
	baseName = strings.TrimPrefix(baseName, "/static/posters/")

	postersDir := viper.GetString("storage.posters-dir")
	fullPath := filepath.Join(postersDir, baseName)

	// Проверяем существование файла
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		logger.Debugf("[usecases] DeletePoster file does not exist: %s", fullPath)
		return nil // Файл не существует - это не ошибка
	}

	// Удаляем файл
	if err := os.Remove(fullPath); err != nil {
		logger.Errorf("[usecases] DeletePoster error removing file: %s", err.Error())
		return err
	}

	logger.Infof("[usecases] DeletePoster successfully deleted file: %s", baseName)
	return nil
}