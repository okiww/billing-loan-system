package errors

import (
	"fmt"
	"net/http"
)

var (
	ErrorInternalServer    = NewError("Internal Server Error", http.StatusInternalServerError)
	ErrorForbiddenResource = NewError("Forbidden Resource", http.StatusForbidden)
	ErrorBadRequest        = NewError("Bad Request", http.StatusBadRequest)
	ErrorNotAcceptable     = NewError("Not Acceptable", http.StatusNotAcceptable)
	ErrorConflict          = NewError("Conflict In The Request", http.StatusConflict)
	ErrorUnauthorized      = NewError("Not Authorized", http.StatusUnauthorized)
	ErrorNotFound          = NewError("Not Found", http.StatusNotFound)
)

type Error struct {
	ErrorMessage string `json:"error,omitempty"`
	ErrorCode    int    `json:"error_code,omitempty"`
}

func (re *Error) Error() string {
	return fmt.Sprintf("%s with code %d", re.ErrorMessage, re.ErrorCode)
}

func NewError(message string, code int) *Error {
	return &Error{
		ErrorMessage: message,
		ErrorCode:    code,
	}
}
