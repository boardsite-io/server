package attachment

type Handler interface {
	Upload(data []byte) (string, error)
	Get(attachID string) ([]byte, string, error)
	Clear() error
}
