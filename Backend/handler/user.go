package handler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kazimovzaman2/Go-chatapp/database"
	"github.com/kazimovzaman2/Go-chatapp/model"
	"github.com/kazimovzaman2/Go-chatapp/utils"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []model.User
	db.Find(&users)
	var responseData []model.UserResponse

	for _, user := range users {
		responseData = append(responseData, utils.UserToResponse(user))
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All users",
		"data":    responseData,
	})
}

func GetMe(c *fiber.Ctx) error {
	user_claim := c.Locals("user").(*jwt.Token)
	claims := user_claim.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db := database.DB
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"errors":  "User not found",
		})
	}

	responseData := utils.UserToResponse(user)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User found",
		"data":    responseData,
	})
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User

	db.Find(&user, id)
	if user.ID == 0 || user.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"errors":  "User not found",
		})
	}

	responseData := utils.UserToResponse(user)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User found",
		"data":    responseData,
	})
}

func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Email        string `json:"email"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		ProfileImage string `json:"profile_image"`
	}

	db := database.DB
	user := new(model.User)

	// Parse request body into user struct
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error(),
		})
	}

	// Validate request body
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	// Check if email already exists
	existingUser := new(model.User)
	if err := db.Where("email = ?", user.Email).First(existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "A user with the provided email already exists",
			"errors":  "Email already exists",
		})
	}

	// Save profile image
	if strings.HasPrefix(user.ProfileImage, "data:@image/") || strings.HasPrefix(user.ProfileImage, "data:@file/") {
		imagePath, err := utils.SaveBase64Image(user.ProfileImage)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error saving image",
				"errors":  err.Error(),
			})
		}

		user.ProfileImage = fmt.Sprintf("http://localhost:8000/%s", imagePath)
	}

	// Hash password
	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't hash password",
			"errors":  err.Error(),
		})
	}
	user.Password = hash

	// Create user
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create user",
			"errors":  err.Error(),
		})
	}

	// Create a new user response
	newUser := NewUser{
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ProfileImage: user.ProfileImage,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Created user",
		"data":    newUser,
	})
}
