package utils

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func NewBadRequest(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}

func NewValidationError(fields map[string]string) ErrorResponse {
	return ErrorResponse{
		Error:  "validation error",
		Fields: fields,
	}
}
