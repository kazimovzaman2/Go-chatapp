package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/kazimovzaman2/Go-chatapp/model"
)

func ValidateUserCredentials(user *model.User) []*model.ErrorResponse {
	var errors []*model.ErrorResponse
	validate := validator.New()

	errs := validate.Struct(user)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			field := err.Field()
			var message string
			switch err.Tag() {
			case "required":
				message = field + " is required"
			case "email":
				message = field + " must be a valid email address"
			case "gte":
				message = field + " must be at least " + err.Param()
			default:
				message = "Validation error on field: " + field
			}
			errors = append(errors, &model.ErrorResponse{
				Status:  "error",
				Message: message,
				Errors:  err.Error(),
			})
		}
	}

	return errors
}
