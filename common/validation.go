package common

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ErrorResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

var errorMessages = map[string]string{
	"alpha":       "%s must contain only alphabetic characters",
	"alphanum":    "%s must contain only alphanumeric characters",
	"email":       "%s must be a valid email address",
	"eq":          "%s must be equal to %s",
	"gt":          "%s must be greater than %s",
	"gte":         "%s must be greater than or equal to %s",
	"hexcolor":    "%s must be a valid hex color",
	"hexadecimal": "%s must be a hexadecimal value",
	"hostname":    "%s must be a valid hostname",
	"ipv4":        "%s must be a valid IPv4 address",
	"ipv6":        "%s must be a valid IPv6 address",
	"isbn10":      "%s must be a valid ISBN-10",
	"isbn13":      "%s must be a valid ISBN-13",
	"lt":          "%s must be less than %s",
	"lte":         "%s must be less than or equal to %s",
	"mac":         "%s must be a valid MAC address",
	"max":         "%s must be at most %s characters long",
	"min":         "%s must be at least %s characters long",
	"ne":          "%s must not be equal to %s",
	"numeric":     "%s must be a numeric value",
	"oneof":       "%s must be one of the following: %s",
	"required":    "%s is required",
	"unique":      "%s must be unique",
	"url":         "%s must be a valid URL",
	"uuid":        "%s must be a valid UUID",
}

func Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errs := validate.Struct(data)

	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			fieldName := err.Field()
			elem.Field = strings.ToLower(fieldName[:1]) + fieldName[1:]
			elem.Error = fmt.Sprintf(errorMessages[err.Tag()], elem.Field, err.Param())
			elem.Error = strings.TrimSuffix(elem.Error, "%!(EXTRA string=)")

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
