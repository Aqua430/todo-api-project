package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	fields := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())

			switch fe.Tag() {
			case "required":
				fields[field] = "is required"
			case "min":
				fields[field] = "must be at least " + fe.Param() + " characters long"
			case "email":
				fields[field] = "is invalid"
			case "lte":
				fields[field] = "must be less than or equal to " + fe.Param()
			case "gte":
				fields[field] = "must be greater than or equal to " + fe.Param()
			default:
				fields[field] = "invalid value"
			}
		}
	}

	return fields
}
