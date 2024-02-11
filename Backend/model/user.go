package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email" form:"email"`
	Password     string `gorm:"not null;" validate:"required,gte=8" json:"password" form:"password"`
	FirstName    string `gorm:"size:255;not null;" validate:"required" json:"first_name" form:"first_name"`
	LastName     string `gorm:"size:255;not null;" validate:"required" json:"last_name" form:"last_name"`
	ProfileImage string `json:"profile_image" form:"profile_image"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ProfileImage string `json:"profile_image"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}
