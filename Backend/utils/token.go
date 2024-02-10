package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kazimovzaman2/Go-chatapp/config"
	"github.com/kazimovzaman2/Go-chatapp/model"
)

func GenerateAccessToken(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Second * 15).Unix()
	accessToken, err := token.SignedString([]byte(config.JWTAccessSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenerateRefreshToken(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	accessToken, err := token.SignedString([]byte(config.JWTRefreshSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
