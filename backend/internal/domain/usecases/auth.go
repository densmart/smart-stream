package usecases

import (
	"net/http"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/utils"
	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/spf13/viper"
)

func SignIn(data dto.SignInRequestDTO) (*dto.AuthResponseDTO, *UCError) {

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
