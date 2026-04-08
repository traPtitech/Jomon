//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/internal/logging"
	"go.uber.org/zap"
)

type FileRepository interface {
	CreateFile(
		ctx context.Context,
		name string,
		mimetype string,
		applicationID uuid.UUID,
		userID uuid.UUID,
	) (*File, error)
	GetFile(ctx context.Context, fileID uuid.UUID) (*File, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
}

type File struct {
	ID        uuid.UUID
	Name      string
	MimeType  string
	CreatedBy uuid.UUID
	CreatedAt time.Time
}

var acceptedMimeTypes = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"image/bmp":          true,
	"application/pdf":    true,
	"application/msword": true,
	"application/zip":    true,
}

func (s *Service) WriteFile(
	ctx context.Context,
	user *User, applicationID uuid.UUID,
	name string, mimetype string, content io.Reader,
) (*File, error) {
	if !acceptedMimeTypes[mimetype] {
		return nil, NewBadInputError("unsupported mime type")
	}
	logger := logging.GetLogger(ctx)
	file, err := s.repository.CreateFile(ctx, name, mimetype, applicationID, user.ID)
	if err != nil {
		logger.Error("failed to create file in repository", zap.Error(err))
		return nil, NewUnexpectedError(err)
	}
	err = s.storage.Save(ctx, file.ID.String(), content)
	if err != nil {
		logger.Error("failed to save file id in storage", zap.Error(err))
		// TODO: storageが返すエラーはそのまま返したい
		return nil, NewUnexpectedError(err)
	}
	return file, nil
}
