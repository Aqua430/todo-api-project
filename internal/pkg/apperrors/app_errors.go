package apperrors

import "net/http"

type AppError struct {
	StatusCode int
	Message    string
	Fields     map[string]string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Message: msg}
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Message: msg}
}

func NewUnauthorizedError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusUnauthorized, Message: msg}
}

func NewValidationError(fields map[string]string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "validation error",
		Fields:     fields,
	}
}
