package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kazimovzaman2/Go-chatapp/config"
	"github.com/kazimovzaman2/Go-chatapp/database"
	"github.com/kazimovzaman2/Go-chatapp/model"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateAccessToken(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Second * 15).Unix()
	accessToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func generateRefreshToken(user model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	accessToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func Login(c *fiber.Ctx) error {
	db := database.DB
	var user model.User

	var input model.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"errors":  err.Error(),
		})
	}

	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"errors":  "User not found",
		})
	}

	if !CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid password",
			"errors":  "Invalid password",
		})
	}

	accessToken, err := generateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not login",
			"errors":  err.Error(),
		})
	}

	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not login",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logged in",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func RefreshToken(c *fiber.Ctx) error {
	type RefreshTokenInput struct {
		RefreshToken string `json:"refresh_token"`
	}

	var input RefreshTokenInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"errors":  err.Error(),
		})
	}

	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
			"errors":  err.Error(),
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
			"errors":  "Invalid token",
		})
	}

	db := database.DB
	userID := claims["id"].(float64)
	var user model.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"errors":  "User not found",
		})
	}

	accessToken, err := generateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not refresh token",
			"errors":  err.Error(),
		})
	}

	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not login",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Token refreshed",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}
