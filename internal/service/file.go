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

var ApprovedMimeTypes = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"image/bmp":          true,
	"application/pdf":    true,
	"application/msword": true,
	"application/zip":    true,
}

func (s *Service) GetFile(ctx context.Context, fileID uuid.UUID) (*File, error) {
	logger := logging.GetLogger(ctx)
	file, err := s.repository.GetFile(ctx, fileID)
	if err != nil {
		logger.Error("failed to get file from repository", zap.Error(err))
		return nil, err
	}
	return file, nil
}

func (s *Service) ReadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, error) {
	logger := logging.GetLogger(ctx)
	content, err := s.storage.Open(ctx, fileID.String())
	if err != nil {
		logger.Error("failed to open file from storage", zap.Error(err))
		return nil, err
	}
	return content, nil
}

func (s *Service) DeleteFile(ctx context.Context, user *User, fileID uuid.UUID) error {
	logger := logging.GetLogger(ctx)
	file, err := s.repository.GetFile(ctx, fileID)
	if err != nil {
		logger.Error("failed to get file from repository", zap.Error(err))
		return err
	}
	if !s.isAccountManagerOrFileCreator(user, file) {
		logger.Info(
			"user is not account manager or file creator",
			zap.String("userID", user.ID.String()),
			zap.String("fileID", fileID.String()))
		return NewForbiddenError("user is not accountManager or file creator")
	}
	// TODO(db-transaction): DeleteFileはストレージとリポジトリ両方で成功させるべき
	err = s.repository.DeleteFile(ctx, fileID)
	if err != nil {
		logger.Error("failed to delete file from repository", zap.Error(err))
		return err
	}
	err = s.storage.Delete(ctx, fileID.String())
	if err != nil {
		logger.Error("failed to delete file from storage", zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) isAccountManagerOrFileCreator(user *User, file *File) bool {
	return s.isAccountManager(user) || file.CreatedBy == user.ID
}

func (s *Service) WriteFile(
	ctx context.Context,
	user *User, applicationID uuid.UUID,
	name string, mimetype string, content io.Reader,
) (*File, error) {
	if !ApprovedMimeTypes[mimetype] {
		return nil, NewBadInputError("unsupported mime type")
	}
	logger := logging.GetLogger(ctx)
	// TODO(db-transaction): WriteFileはストレージとリポジトリ両方で成功させるべき
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
