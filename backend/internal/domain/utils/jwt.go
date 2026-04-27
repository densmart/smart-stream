package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	accessTokenTTL = 24 * time.Hour
)

type JwtAuth struct {
	apiKey string
	name   string
}

func NewJwtAuth(apiKey string, userName string) *JwtAuth {
	return &JwtAuth{apiKey: apiKey, name: userName}
}

type JwtToken struct {
	apiKey string
	Access string
}

func NewJwtToken(apiKey string) *JwtToken {
	return &JwtToken{apiKey: apiKey}
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserName string `json:"user_name"`
}

// GenerateAccessToken generates access token
func (ja *JwtAuth) GenerateAccessToken() (string, error) {
	claims := AccessTokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(accessTokenTTL),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		ja.name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(ja.apiKey))
}

// GetAccessClaims checks access authentication and returns access token claims
func (tc *JwtToken) GetAccessClaims() (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tc.Access, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(tc.apiKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, errors.New("invalid access token claims")
	}
	return claims, nil
}
