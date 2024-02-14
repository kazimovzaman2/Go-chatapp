package utils

import "github.com/kazimovzaman2/Go-chatapp/model"

func UserToResponse(user model.User) model.UserResponse {
	return model.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
