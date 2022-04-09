package customErrors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var Trans ut.Translator

// ParseError - parses error and returns appropriate views.ErrorResponse.
func ParseError(err error) *ErrorResponse {

	t := reflect.TypeOf(err)
	fmt.Println(t)

	switch t := err.(type) {

	case validator.ValidationErrors:
		return handleValidationErr(t)

	case *json.UnmarshalTypeError:
		return handleUnMarshalTypeError(t)

	case *json.SyntaxError:
		// fmt.Println("invalid unmarshal error")
		return New(http.StatusBadRequest, err)

	case *ErrorResponse:
		return t

	default:

		switch {
		case errors.Is(err, io.EOF):
			return New(http.StatusBadRequest, err)
		case errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), sql.ErrNoRows.Error()):
			return New(http.StatusNotFound, NotFound)

		default:
			return NewInternalServerError()

		}
	}
}

// handleUnMarshalTypeError - handles errors returned by json unmarshalling and returns appropriate views.ErrorResponse
func handleUnMarshalTypeError(err *json.UnmarshalTypeError) (resp *ErrorResponse) {
	resp = &ErrorResponse{
		ErrStatus: http.StatusBadRequest,
	}

	resp.Errors = append(resp.Errors, ErrorValidation{
		Field:   err.Field,
		Message: err.Error(),
		Type:    err.Type.String(),
	})

	return
}

// handleValidationErr - handles errors returned by go validator package and returns appropriate views.ErrorResponse
func handleValidationErr(err validator.ValidationErrors) (resp *ErrorResponse) {
	resp = &ErrorResponse{
		ErrStatus: http.StatusBadRequest,
	}

	for _, v := range err {

		errVal := ErrorValidation{
			Field:   v.Field(),
			Message: v.Translate(Trans),
		}

		resp.Errors = append(resp.Errors, errVal)
	}

	return
}
