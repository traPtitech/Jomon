package service

import (
	"fmt"
	"io"
	"mime"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        uuid.UUID
	Name      string
	MimeType  string
	CreatedAt time.Time
}

func (s *Services) CreateFile(src io.Reader, id uuid.UUID, mimetype string) error {
	ext, err := mime.ExtensionsByType(mimetype)
	if err != nil {
		return err
	} else if len(ext) == 0 {
		return fmt.Errorf("%s is not registered", mimetype)
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext[0])

	err = s.Storage.Save(filename, src)

	return err
}

func (s *Services) OpenFile(fileID uuid.UUID, mimetype string) (io.ReadCloser, error) {
	ext, err := mime.ExtensionsByType(mimetype)
	if err != nil {
		return nil, err
	} else if len(ext) == 0 {
		return nil, fmt.Errorf("%s is not registered", mimetype)
	}

	filename := fmt.Sprintf("%s%s", fileID.String(), ext[0])

	return s.Storage.Open(filename)
}

func (s *Services) DeleteFile(fileID uuid.UUID, mimetype string) error {
	ext, err := mime.ExtensionsByType(mimetype)
	if err != nil {
		return err
	} else if len(ext) == 0 {
		return fmt.Errorf("%s is not registered", mimetype)
	}
	filename := fmt.Sprintf("%s%s", fileID.String(), ext[0])

	return s.Storage.Delete(filename)
}
