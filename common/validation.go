package common

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errs := validate.Struct(data)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.Field = err.Field()
			elem.Error = err.Tag()

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
