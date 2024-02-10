package handler

import (
	"fmt"
	"os"
	"path/filepath"

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

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "All users",
		Data:    responseData,
	})
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User

	db.Find(&user, id)
	if user.ID == 0 || user.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "User with the provided ID not found",
			Errors:  "User not found",
		})
	}

	responseData := utils.UserToResponse(user)

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "User found",
		Data:    responseData,
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
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Couldn't parse request body",
			Errors:  err.Error(),
		})
	}

	// Validate request body
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Invalid input",
			Errors:  err.Error(),
		})
	}

	// Check if email already exists
	existingUser := new(model.User)
	if err := db.Where("email = ?", user.Email).First(existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "User with this email already exists",
			Errors:  "Email already exists",
		})
	}

	// Save profile image
	if utils.IsBase64(user.ProfileImage) {
		imagePath, err := utils.SaveBase64Image(user.ProfileImage)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Status:  "error",
				Message: "Couldn't save profile image",
				Errors:  err.Error(),
			})
		}

		user.ProfileImage = fmt.Sprintf("http://localhost:8000/%s", imagePath)
	}

	// Hash password
	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Couldn't hash password",
			Errors:  err.Error(),
		})
	}
	user.Password = hash

	// Create user
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Couldn't create user",
			Errors:  err.Error(),
		})
	}

	accessToken, refreshToken, err := generateTokens(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Could not login",
			Errors:  err.Error(),
		})
	}

	// Create a new user response
	newUser := NewUser{
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ProfileImage: user.ProfileImage,
	}

	return c.Status(fiber.StatusCreated).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "User created",
		Data: fiber.Map{
			"user":          newUser,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

func GetMe(c *fiber.Ctx) error {
	user_claim := c.Locals("user").(*jwt.Token)
	claims := user_claim.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db := database.DB
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "User not found",
			Errors:  err.Error(),
		})
	}

	responseData := utils.UserToResponse(user)

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "User found",
		Data:    responseData,
	})
}

func UpdateMe(c *fiber.Ctx) error {
	user_claim := c.Locals("user").(*jwt.Token)
	claims := user_claim.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db := database.DB
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "User not found",
			Errors:  err.Error(),
		})
	}

	// Parse request body into user struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Couldn't parse request body",
			Errors:  err.Error(),
		})
	}

	// Save profile image
	if utils.IsBase64(user.ProfileImage) {
		imagePath, err := utils.SaveBase64Image(user.ProfileImage)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Status:  "error",
				Message: "Couldn't save profile image",
				Errors:  err.Error(),
			})
		}

		user.ProfileImage = fmt.Sprintf("http://localhost:8000/%s", imagePath)
	}

	// Update user
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Couldn't update user",
			Errors:  err.Error(),
		})
	}

	responseData := utils.UserToResponse(user)

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "User updated",
		Data:    responseData,
	})
}

func DeleteMe(c *fiber.Ctx) error {
	user_claim := c.Locals("user").(*jwt.Token)
	claims := user_claim.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db := database.DB
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "User not found",
			Errors:  err.Error(),
		})
	}

	imagePath := user.ProfileImage

	if err := db.Unscoped().Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Status:  "error",
			Message: "Couldn't delete user",
			Errors:  err.Error(),
		})
	}

	if imagePath != "" {
		filename := filepath.Base(imagePath)
		err := os.Remove(filepath.Join("./media/avatars", filename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Status:  "error",
				Message: "Couldn't delete profile image",
				Errors:  err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Status:  "success",
		Message: "User deleted",
		Data:    nil,
	})
}
