package attachment

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/h2non/filetype"
	gonanoid "github.com/matoous/go-nanoid"
)

const attachmentDir = "/tmp/attachment"

var (
	ErrCreate   = errors.New("cannot allocate resource")
	ErrWrite    = errors.New("cannot write file content")
	ErrFileType = errors.New("unsupported MIME type")
	ErrNotFound = errors.New("file not found")
)

type localAttachment struct {
	baseDir string
}

// NewLocalHandler create a new attachment Handler for storing attachment in the local filesystem.
func NewLocalHandler(sessionID string) Handler {
	return &localAttachment{
		baseDir: fmt.Sprintf("%s/%s", attachmentDir, sessionID),
	}
}

func (a *localAttachment) Upload(data []byte) (string, error) {
	fileExt, err := getFileExtension(data)
	if err != nil {
		return "", err
	}
	attachID := fmt.Sprintf("%s.%s", gonanoid.MustID(32), fileExt)

	if err := os.MkdirAll(a.baseDir, 0666); err != nil {
		return "", ErrCreate
	}
	file, err := os.OpenFile(
		fmt.Sprintf("%s/%s", a.baseDir, attachID),
		os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return "", ErrCreate
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return "", ErrWrite
	}
	return attachID, nil
}

func (a *localAttachment) Get(attachID string) (io.Reader, string, error) {
	f, err := os.Open(fmt.Sprintf("%s/%s", a.baseDir, attachID))
	if err != nil {
		return nil, "", ErrNotFound
	}
	fType, err := filetype.MatchReader(f)
	if err != nil {
		return nil, "", ErrNotFound
	}
	_, err = f.Seek(0, 0)
	return f, fType.MIME.Value, err
}

func (a *localAttachment) Clear() error {
	return os.RemoveAll(a.baseDir)
}

func getFileExtension(data []byte) (string, error) {
	fType, err := filetype.Match(data)
	if err != nil {
		return "", ErrFileType
	}
	if fType.MIME.Value != "application/pdf" && fType.MIME.Value != "image/png" {
		return "", ErrFileType
	}
	return fType.Extension, nil
}
