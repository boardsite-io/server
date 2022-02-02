package attachment

import "io"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . Handler
type Handler interface {
	Upload(data []byte) (string, error)
	Get(attachID string) (io.Reader, string, error)
	Clear() error
}

type AttachmentResponse struct {
	AttachID string `json:"attachId"`
}
