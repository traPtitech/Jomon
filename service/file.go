package service

import (
	"io"

	"github.com/google/uuid"
)

type File struct {
	ID       uuid.UUID
	MimeType string
}

func (s *Services) CreateFile(src io.Reader, mimetype string) (File, error) {
	// TODO: Implement
	return File{}, nil
}
