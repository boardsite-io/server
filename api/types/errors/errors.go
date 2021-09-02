package errors

import (
	"fmt"
	"net/http"

	"github.com/heat1q/boardsite/api/types"
)

var (
	BadRequest = NewWrapper(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	NotFound   = NewWrapper(http.StatusNotFound, http.StatusText(http.StatusNotFound))

	InternalServerError = NewWrapper(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
)

// Wrapper defines a generic error wrapper
type Wrapper struct {
	Code int    `json:"-"`
	Info string `json:"info"`
	types.Message
}

// New returns a generic error that map to http.StatusInternalServerError by default
func New(info string) error {
	return &Wrapper{
		Message: types.Message{
			Type: "error",
		},
		Code: http.StatusInternalServerError,
		Info: info,
	}
}

// NewWrapper create a new Wrapper with any valid HTTP status code
func NewWrapper(code int, info string) *Wrapper {
	return &Wrapper{
		Message: types.Message{
			Type: "error",
		},
		Code: code,
		Info: info,
	}
}

func (e *Wrapper) Error() string {
	return e.Info
}

// SetInfo sets the error information
func (e *Wrapper) SetInfo(v interface{}) *Wrapper {
	e.Info = fmt.Sprintf("%v", v)
	return e
}

// MaptoHTTPError maps any error to an suitable error with HTTP status code
func MaptoHTTPError(err error) error {
	wrapper, ok := err.(*Wrapper)
	if !ok {
		return InternalServerError.SetInfo(err)
	}
	return wrapper
}
