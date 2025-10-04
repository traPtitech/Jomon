//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
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
