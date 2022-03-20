package errors

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrBadRequest          = New(http.StatusBadRequest)
	ErrUnauthorized        = New(http.StatusUnauthorized)
	ErrForbidden           = New(http.StatusForbidden)
	ErrNotFound            = New(http.StatusNotFound)
	ErrInternalServerError = New(http.StatusInternalServerError)
	ErrBadGateway          = New(http.StatusBadGateway)
)

type HTTPError struct {
	Status   int    `json:"-"`
	Message  string `json:"message"`
	Code     Code   `json:"code,omitempty"`
	internal error
}

var _ error = (*HTTPError)(nil)

func New(status int, message ...string) *HTTPError {
	msg := http.StatusText(status)
	if len(message) > 0 {
		msg = message[0]
	}
	return &HTTPError{
		Status:  status,
		Message: msg,
	}
}

// From creates a new HTTPError from a specific error code.
func From(code Code) *HTTPError {
	status, ok := codeStatusMap[code]
	if !ok {
		status = http.StatusInternalServerError
	}
	httpErr := New(status)
	httpErr.Code = code
	return httpErr
}

type ErrorOption func(e *HTTPError)

func WithError(err error) ErrorOption {
	return func(e *HTTPError) {
		e.internal = err
	}
}

func WithErrorf(format string, a ...interface{}) ErrorOption {
	return WithError(fmt.Errorf(format, a...))
}

func (e *HTTPError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("status=%d code=%d", e.Status, e.Code))
	if e.Message != "" {
		sb.WriteString(fmt.Sprintf(": %s", e.Message))
	}
	if e.internal != nil {
		sb.WriteString(fmt.Sprintf(": %v", e.internal))
	}
	return sb.String()
}

func (e HTTPError) Is(err error) bool {
	target, ok := err.(*HTTPError)
	if !ok {
		return false
	}
	return e.Status == target.Status && e.Code == target.Code
}

func (e HTTPError) Wrap(options ...ErrorOption) error {
	for _, o := range options {
		o(&e)
	}
	return &e
}

func (e *HTTPError) Unwrap() error {
	return e.internal
}
