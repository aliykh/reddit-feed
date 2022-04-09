package customErrors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorResponse for error responses.
type ErrorResponse struct {
	ErrStatus int               `json:"status,omitempty"`
	ErrError  string            `json:"error,omitempty"`
	Errors    []ErrorValidation `json:"errors,omitempty"`
}

// New returns a new views.ErrorResponse
func New(code int, err error) *ErrorResponse {
	return &ErrorResponse{
		ErrStatus: code,
		ErrError:  err.Error(),
	}
}

// ErrorValidation for validation errors.
type ErrorValidation struct {
	Field   string `json:"field"`
	Message string `json:"message,omitempty"`
	Type    string `json:"type,omitempty"`
}

// ErrorResponse implements error interface.
func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("status: %d; error: %v;\n", err.ErrStatus, err.ErrError)
}

// NewInternalServerError - returns a new views.ErrorResponse (status internal server error 500)
func NewInternalServerError() *ErrorResponse {
	return &ErrorResponse{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
	}
}

var (
	NotFound              = errors.New("not found")
	InternalServerError   = errors.New("internal server error")
	InvalidUriParam       = errors.New("invalid uri param")
)
