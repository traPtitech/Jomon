package model

import (
	"fmt"
	"github.com/gofrs/uuid"
	storagePkg "github.com/traPtitech/Jomon/storage"
	"io"
	"mime"
)

type ApplicationsImage struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key"`
	ApplicationID uuid.UUID `gorm:"type:char(36);not null"`
	MimeType      string    `gorm:"type:text;not null" json:"-"`
}

type ApplicationsImageRepository interface {
	CreateApplicationsImage(applicationId uuid.UUID, src io.Reader, mimeType string) (uuid.UUID, error)
	OpenApplicationsImage(id uuid.UUID) (io.ReadCloser, error)
	DeleteApplicationsImage(id uuid.UUID) error
}

type applicationsImageRepository struct {
	storage storagePkg.Storage
}

func NewApplicationsImageRepository(storage storagePkg.Storage) ApplicationsImageRepository {
	return &applicationsImageRepository{storage: storage}
}

func (repo *applicationsImageRepository) CreateApplicationsImage(applicationId uuid.UUID, src io.Reader, mimeType string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return uuid.Nil, err
	} else if len(ext) == 0 {
		return uuid.Nil, fmt.Errorf("%s is not registered", mimeType)
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext[0])

	err = repo.storage.Save(filename, src)
	if err != nil {
		return uuid.Nil, err
	}

	im := ApplicationsImage{
		ID:            id,
		ApplicationID: applicationId,
		MimeType:      mimeType,
	}

	if err = db.Create(&im).Error; err != nil {
		_ = repo.storage.Delete(filename)
		return uuid.Nil, err
	}

	return id, nil
}

func (repo *applicationsImageRepository) OpenApplicationsImage(id uuid.UUID) (io.ReadCloser, error) {
	im := ApplicationsImage{
		ID: id,
	}

	if err := db.First(&im).Error; err != nil {
		return nil, err
	}

	ext, err := mime.ExtensionsByType(im.MimeType)
	if err != nil {
		return nil, err
	} else if len(ext) == 0 {
		return nil, fmt.Errorf("%s is not registered", im.MimeType)
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext[0])

	return repo.storage.Open(filename)
}

func (repo *applicationsImageRepository) DeleteApplicationsImage(id uuid.UUID) error {
	im := ApplicationsImage{
		ID: id,
	}

	if err := db.First(&im).Error; err != nil {
		return err
	}

	ext, err := mime.ExtensionsByType(im.MimeType)
	if err != nil {
		return err
	} else if len(ext) == 0 {
		return fmt.Errorf("%s is not registered", im.MimeType)
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext[0])

	if err := db.Delete(im).Error; err != nil {
		return err
	}

	return repo.storage.Delete(filename)
}
