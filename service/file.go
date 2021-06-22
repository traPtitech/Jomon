package service

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/model"
)

type File struct {
	ID        uuid.UUID
	Name      string
	MimeType  string
	CreatedAt time.Time
}

func (s *Services) CreateFile(src io.Reader, name string, mimetype string, requestID uuid.UUID) (*File, error) {
	ctx := context.Background()
	modelfile, err := s.Repository.CreateFile(ctx, src, name, mimetype, requestID)
	if err != nil {
		return nil, err
	}

	return ConvertModelFileToServiceFile(modelfile), nil
}

func (s *Services) GetFile(fileID uuid.UUID) (*File, error) {
	ctx := context.Background()
	file, err := s.Repository.GetFile(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return ConvertModelFileToServiceFile(file), nil
}

func (s *Services) OpenFile(fileID uuid.UUID) (io.ReadCloser, error) {
	ctx := context.Background()
	f, err := s.Repository.OpenFile(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Services) DeleteFile(fileID uuid.UUID, requestID uuid.UUID) error {
	ctx := context.Background()
	return s.Repository.DeleteFile(ctx, fileID, requestID)
}

func ConvertModelFileToServiceFile(modelfile *model.File) *File {
	return &File{
		ID:        modelfile.ID,
		Name:      modelfile.Name,
		MimeType:  modelfile.MimeType,
		CreatedAt: modelfile.CreatedAt,
	}
}
