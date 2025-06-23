package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	localDir string
}

var _ Storage = (*Local)(nil)

func NewLocalStorage(dir string) (*Local, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return &Local{}, errors.New("dir doesn't exist")
	}
	if !fi.IsDir() {
		return &Local{}, errors.New("dir is not a directory")
	}

	return &Local{localDir: dir}, nil
}

func (l *Local) Save(ctx context.Context, filename string, src io.Reader) error {
	file, err := os.Create(l.getFilePath(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}

func (l *Local) Open(ctx context.Context, filename string) (io.ReadCloser, error) {
	r, err := os.Open(l.getFilePath(filename))
	if err != nil {
		return nil, ErrFileNotFound
	}
	return r, nil
}

func (l *Local) Delete(ctx context.Context, filename string) error {
	path := l.getFilePath(filename)
	if _, err := os.Stat(path); err != nil {
		return ErrFileNotFound
	}
	return os.Remove(path)
}

func (l *Local) getFilePath(filename string) string {
	return filepath.Join(l.localDir, filename)
}
