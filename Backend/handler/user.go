package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kazimovzaman2/Go-chatapp/database"
	"github.com/kazimovzaman2/Go-chatapp/model"
	"github.com/kazimovzaman2/Go-chatapp/utils"
	"github.com/kazimovzaman2/Go-chatapp/validation"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GetAllUsers is a handler to get all users
// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Router /users/ [get]
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

// GetUser is a handler to get a user by ID
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /users/{id} [get]
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

// CreateUser is a handler to create a new user
// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.User true "Create user"
// @Success 201 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/ [post]
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

	userValidationErrors := validation.ValidateUserCredentials(user)
	if len(userValidationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(userValidationErrors)
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

// GetMe is a handler to get the current user
// @Summary Get the current user
// @Description Get the current user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /users/me/ [get]
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

// UpdateMe is a handler to update the current user
// @Summary Update the current user
// @Description Update the current user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body model.User true "User data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/ [patch]
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

// DeleteMe is a handler to delete the current user
// @Summary Delete the current user
// @Description Delete the current user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/me/ [delete]
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
