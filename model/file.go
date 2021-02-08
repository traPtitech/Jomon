package model

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"time"

	"github.com/gofrs/uuid"
	storagePkg "github.com/traPtitech/Jomon/storage"
)

// File struct fo File
type File struct {
	ID        uuid.UUID  `gorm:"type:char(36);primary_key"`
	RequestID uuid.UUID  `gorm:"type:char(36);not null"`
	MimeType  string     `gorm:"type:text;not null" json:"-"`
	CreatedAt time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"deleted_at"`
}

// MarshalJSON Marshal File struct
func (file File) MarshalJSON() ([]byte, error) {
	return json.Marshal(file.ID)
}

// FileRepository Repo of File
type FileRepository interface {
	CreateFile(requestID uuid.UUID, src io.Reader, mimeType string) (File, error)
	GetFile(id uuid.UUID) (File, error)
	OpenFile(file File) (io.ReadCloser, error)
	DeleteFile(file File) error
}

type fileRepository struct {
	storage storagePkg.Storage
}

// NewFileRepository Make FileRepository
func NewFileRepository(storage storagePkg.Storage) FileRepository {
	return &fileRepository{storage: storage}
}

func (repo *fileRepository) CreateFile(requestID uuid.UUID, src io.Reader, mimeType string) (File, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return File{}, err
	}

	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return File{}, err
	} else if len(ext) == 0 {
		return File{}, fmt.Errorf("%s is not registered", mimeType)
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext[0])

	err = repo.storage.Save(filename, src)
	if err != nil {
		return File{}, err
	}

	im := File{
		ID:        id,
		RequestID: requestID,
		MimeType:  mimeType,
	}

	if err = db.Create(&im).Error; err != nil {
		_ = repo.storage.Delete(filename)
		return File{}, err
	}

	return im, nil
}

func (repo *fileRepository) GetFile(id uuid.UUID) (File, error) {
	im := File{
		ID: id,
	}

	if err := db.First(&im).Error; err != nil {
		return File{}, err
	}

	return im, nil
}

func (repo *fileRepository) OpenFile(file File) (io.ReadCloser, error) {
	ext, err := mime.ExtensionsByType(file.MimeType)
	if err != nil {
		return nil, err
	} else if len(ext) == 0 {
		return nil, fmt.Errorf("%s is not registered", file.MimeType)
	}

	filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

	return repo.storage.Open(filename)
}

func (repo *fileRepository) DeleteFile(file File) error {
	ext, err := mime.ExtensionsByType(file.MimeType)
	if err != nil {
		return err
	} else if len(ext) == 0 {
		return fmt.Errorf("%s is not registered", file.MimeType)
	}

	filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

	if err := db.Delete(file).Error; err != nil {
		return err
	}

	return repo.storage.Delete(filename)
}
