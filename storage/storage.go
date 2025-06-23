//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package storage

import (
	"context"
	"errors"
	"io"
)

var (
	ErrFileNotFound = errors.New("not found")
)

type Storage interface {
	Save(ctx context.Context, filename string, src io.Reader) error
	Open(ctx context.Context, filename string) (io.ReadCloser, error)
	Delete(ctx context.Context, filename string) error
}
