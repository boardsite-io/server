package session

import (
	"errors"
	"io"
	"os"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func UploadAttachment(scb *ControlBlock, src io.Reader, name string) (string, error) {
	if err := os.MkdirAll("./attachments/"+scb.ID, 0666); err != nil {
		return "", errors.New("cannot create resource")
	}
	attachID := gonanoid.MustGenerate(alphabet, 16) + name
	file, err := os.OpenFile(
		"./attachments/"+scb.ID+"/"+attachID,
		os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return "", errors.New("cannot create resource")
	}
	defer file.Close()

	if _, err := io.Copy(file, src); err != nil {
		return "", errors.New("cannot write content")
	}
	return attachID, nil
}

func OpenAttachment(scb *ControlBlock, attachID string) (io.ReadCloser, error) {
	return os.Open("./attachments/" + scb.ID + "/" + attachID)
}

func ClearAttachments(scb *ControlBlock) error {
	return os.RemoveAll("./attachments/" + scb.ID)
}
