package tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//MalformedRequest Error requerido
type MalformedRequest struct {
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}
//ErrorHandling errors
func ErrorHandling(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		Msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case errors.Is(err, io.ErrUnexpectedEOF):
		Msg := fmt.Sprintf("Request body contains badly-formed JSON")
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case errors.As(err, &unmarshalTypeError):
		Msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		Msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case errors.Is(err, io.EOF):
		Msg := "Request body must not be empty"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case err.Error() == "http: request body too large":
		Msg := "Request body must not be larger than 10MB"
		return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: Msg}

	case strings.HasPrefix(err.Error(), "Error parsing"):
		Msg := fmt.Sprintf(err.Error())
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}
	default:
		return err
	}
}
