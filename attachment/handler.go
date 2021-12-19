package attachment

import "io"

type Handler interface {
	Upload(data []byte) (string, error)
	Get(attachID string) (io.Reader, string, error)
	Clear() error
}
