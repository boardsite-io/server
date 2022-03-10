package errors

import "net/http"

type Code uint32

// Client error codes
const (
	BadRequest Code = iota + 4000
	RateLimitExceeded
	MissingIdentifier
	AttachmentSizeExceeded
	MaxNumberOfUsersReached
	BadUsername
	WrongPassword
)

// Server error codes
const (
	CodeInternalError Code = iota + 5000
)

var codeStatusMap = map[Code]int{
	RateLimitExceeded:       http.StatusTooManyRequests,
	MissingIdentifier:       http.StatusForbidden,
	AttachmentSizeExceeded:  http.StatusBadRequest,
	MaxNumberOfUsersReached: http.StatusBadRequest,
	BadUsername:             http.StatusBadRequest,
	WrongPassword:           http.StatusBadRequest,
}
