package usecases

import (
	"net/http"
	"time"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/densmart/smart-stream/internal/domain/utils"
	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/spf13/viper"
)

// WebSignIn - авторизация администратора (для WebAPI)
func WebSignIn(data dto.SignInRequestDTO) (*dto.AuthResponseDTO, *UCError) {

	configLogin := viper.GetString("admin-user.login")
	configPassword := viper.GetString("admin-user.password")

	if data.Login != configLogin || data.Password != configPassword {
		return nil, &UCError{
			Code:     100,
			HttpCode: http.StatusUnauthorized,
		}
	}

	// generate JWT
	accessToken, ucErr := generateJWT(data.Login)
	if ucErr != nil {
		return nil, ucErr
	}

	return &dto.AuthResponseDTO{
		AccessToken: accessToken,
	}, nil
}

// ClientSignIn - авторизация клиента (для ClientAPI)
func ClientSignIn(oltp repo.OltpRepo, data dto.SignInRequestDTO) (*dto.AuthResponseDTO, *UCError) {
	// Поиск клиента по логину
	client, dbErr := oltp.RetrieveClientByLogin(data.Login)
	if dbErr != nil {
		return nil, &UCError{
			Code:     101,
			HttpCode: http.StatusUnauthorized,
			Message:  "Invalid credentials",
		}
	}

	// Проверка is_active флага
	if !client.IsActive {
		return nil, &UCError{
			Code:     102,
			HttpCode: http.StatusForbidden,
			Message:  "Client is not active",
		}
	}

	// Проверка пароля
	if !utils.CheckPasswordHash(data.Password, client.Password) {
		return nil, &UCError{
			Code:     101,
			HttpCode: http.StatusUnauthorized,
			Message:  "Invalid credentials",
		}
	}

	// Обновление last_login_at
	now := time.Now()
	updateData := dto.UpdateClientDTO{
		LastLoginAt: &now,
	}
	_, dbErr = oltp.UpdateClient(client.ID, updateData)
	if dbErr != nil {
		// Логируем ошибку, но не прерываем авторизацию
		logger.Debugf("failed to update last_login_at for client %s: %v", client.ID, dbErr)
	}

	// Генерация JWT
	accessToken, ucErr := generateJWT(client.Login)
	if ucErr != nil {
		return nil, ucErr
	}

	return &dto.AuthResponseDTO{
		AccessToken: accessToken,
	}, nil
}

func generateJWT(login string) (string, *UCError) {
	apiKey := viper.GetString("app.jwt-secret")

	jwtWrapper := utils.NewJwtAuth(apiKey, login)
	accessToken, err := jwtWrapper.GenerateAccessToken()
	if err != nil {
		logger.Debugf("error generating access token: %s", err.Error())
		return "", &UCError{
			Code:     100,
			HttpCode: http.StatusUnauthorized,
		}
	}

	return accessToken, nil
}
