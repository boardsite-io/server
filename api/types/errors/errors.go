package errors

import (
	"fmt"
	"net/http"

	"github.com/heat1q/boardsite/api/types"
)

var (
	BadRequest = NewWrapper(1000, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	NotFound   = NewWrapper(2000, http.StatusNotFound, http.StatusText(http.StatusNotFound))

	InternalServerError = NewWrapper(3000, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
)

// Wrapper defines a generic error wrapper
type Wrapper struct {
	// http status code
	StatusCode int `json:"-"`
	// generic text for the user
	Text string `json:"text,omitempty"`
	// code used by the FE to identify the error
	Code int `json:"code,omitempty"`
	// internal info not transmitted
	Info string `json:"-"`
	types.Message
}

// New returns a generic error that maps to http.StatusInternalServerError by default
func New(info string) error {
	return &Wrapper{
		Message: types.Message{
			Type: "error",
		},
		Info: info,
	}
}

// NewWrapper create a new Wrapper with any valid HTTP status code
func NewWrapper(code, status int, text string) *Wrapper {
	return &Wrapper{
		Message: types.Message{
			Type: "error",
		},
		Code:       code,
		StatusCode: status,
		Text:       text,
	}
}

func (e *Wrapper) Error() string {
	return e.Info
}

// SetInfo sets the error internal information
func (e *Wrapper) SetInfo(v interface{}) *Wrapper {
	e.Info = fmt.Sprintf("%v", v)
	return e
}

// MaptoHTTPError maps any error to an suitable error with HTTP status code
func MaptoHTTPError(err error) error {
	wrapper, ok := err.(*Wrapper)
	if !ok || wrapper.StatusCode == 0 {
		return InternalServerError.SetInfo(err)
	}
	return wrapper
}
