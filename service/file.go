package service

import (
	"context"
	"fmt"
	"io"
	"mime"
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

func (s *Services) DeleteFile(fileID uuid.UUID) error {
	ctx := context.Background()
	return s.Repository.DeleteFile(ctx, fileID)
}

func ConvertModelFileToServiceFile(modelfile *model.File) *File {
	return &File{
		ID:        modelfile.ID,
		Name:      modelfile.Name,
		MimeType:  modelfile.MimeType,
		CreatedAt: modelfile.CreatedAt,
	}
}
