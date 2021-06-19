package service

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type File struct {
	ID       uuid.UUID
	Name     string
	MimeType string
}

func (s *Services) CreateFile(src io.Reader, name string, mimetype string, requestID uuid.UUID) (*File, error) {
	ctx := context.Background()
	modelfile, err := s.Repository.CreateFile(ctx, src, name, mimetype, requestID)
	if err != nil {
		return nil, err
	}

	file := &File{
		ID:       modelfile.ID,
		Name:     modelfile.Name,
		MimeType: modelfile.MimeType,
	}

	return file, nil
}
