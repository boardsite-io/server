package errors

import "net/http"

// Client error codes
const (
	CodeBadRequest uint32 = iota + 4000
	CodeRateLimitExceeded
	CodeMissingIdentifier
	CodeAttachmentSizeExceeded
)

// Server error codes
const (
	CodeInternalError uint32 = iota + 5000
)

var codeStatusMap = map[uint32]int{
	CodeRateLimitExceeded:      http.StatusTooManyRequests,
	CodeMissingIdentifier:      http.StatusForbidden,
	CodeAttachmentSizeExceeded: http.StatusBadRequest,
}
