package errors

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrBadRequest   = New(http.StatusBadRequest)
	ErrUnauthorized = New(http.StatusUnauthorized)
	ErrForbidden    = New(http.StatusForbidden)
	ErrNotFound     = New(http.StatusNotFound)

	ErrInternalServerError = New(http.StatusInternalServerError)
	ErrBadGateway          = New(http.StatusBadGateway)
	ErrGatewayTimeout      = New(http.StatusGatewayTimeout)
)

type HTTPError struct {
	Status   int    `json:"-"`
	Message  string `json:"message"`
	internal error
	Code     Code `json:"code,omitempty"`
}

var _ error = (*HTTPError)(nil)

// New creates a new HTTPError error with the given status and optional message.
// If no message is provided, the default http.StatusText is set as message.
func New(status int, message ...string) *HTTPError {
	msg := http.StatusText(status)
	if len(message) > 0 {
		msg = strings.Join(message, ", ")
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

// WithMessage overwrites the default http.StatusText in the message field.
// This functional argument is passed to the HTTPError.Wrap method.
func WithMessage(msg string) ErrorOption {
	return func(e *HTTPError) {
		e.Message = msg
	}
}

// WithMessagef formats and overwrites the default http.StatusText in the message field.
// This functional argument is passed to the HTTPError.Wrap method.
func WithMessagef(format string, a ...interface{}) ErrorOption {
	return WithMessage(fmt.Sprintf(format, a...))
}

// WithError sets the reason of the error, used for internal logging.
// This functional argument is passed to the HTTPError.Wrap method.
func WithError(err error) ErrorOption {
	return func(e *HTTPError) {
		e.internal = err
	}
}

// WithCode sets the error code.
// This functional argument is passed to the HTTPError.Wrap method.
func WithCode(code Code) ErrorOption {
	return func(e *HTTPError) {
		e.Code = code
	}
}

// WithErrorf formats an internal error.
// This functional argument is passed to the HTTPError.Wrap method.
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

// Wrap wraps any optional argument in the HTTPError
func (e HTTPError) Wrap(options ...ErrorOption) error {
	for _, o := range options {
		o(&e)
	}
	return &e
}

// Unwrap unwraps the internal error of HTTPError.
func (e *HTTPError) Unwrap() error {
	return e.internal
}
