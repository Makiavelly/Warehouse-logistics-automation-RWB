package errors

import "fmt"

type ErrInvalidAuthToken struct {
	Message string `json:"error"`
}

func NewErrInvalidAuthToken(msg string) *ErrInvalidAuthToken {
	return &ErrInvalidAuthToken{Message: msg}
}

func (e *ErrInvalidAuthToken) Error() string {
	return fmt.Sprintf("invalid auth token: %s", e.Message)
}

type ErrResponse struct {
	Error string `json:"error"`
}

func NewErrResponse(msg string) ErrResponse {
	return ErrResponse{Error: msg}
}