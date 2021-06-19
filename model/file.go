package model

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type FileRepository interface {
	CreateFile(ctx context.Context, src io.Reader, name string, mimetype string, requestID uuid.UUID) (*File, error)
}

type File struct {
	ID        uuid.UUID
	Name      string
	MimeType  string
	CreatedAt time.Time
}
