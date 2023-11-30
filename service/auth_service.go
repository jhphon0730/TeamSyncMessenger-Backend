package service

import (
	"TeamSyncMessenger-Backend/DTO"
	"TeamSyncMessenger-Backend/model"
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

type AuthService interface {
	CreateUserLoginJWT(user DTO.LoginUserDTO) (string, error)
}

type authService struct {
}

func NewAuthService() *authService {
	return &authService{}
}

func (ac *authService) CreateUserLoginJWT(user DTO.LoginUserDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	tokenStr, err := token.SignedString(model.JwtKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil

}
