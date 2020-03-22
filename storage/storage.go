package storage

import (
	"errors"
	"io"
)

var (
	ErrFileNotFound = errors.New("not found")
)

type Storage interface {
	Save(filename string, src io.Reader) error
	Open(filename string) (io.ReadCloser, error)
	Delete(filename string) error
}
