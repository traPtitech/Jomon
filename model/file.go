//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FileRepository interface {
	CreateFile(ctx context.Context, name string, mimetype string, requestID uuid.UUID) (*File, error)
	GetFile(ctx context.Context, fileID uuid.UUID) (*File, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
}

type File struct {
	ID        uuid.UUID
	Name      string
	MimeType  string
	CreatedAt time.Time
}
