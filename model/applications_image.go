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

type ApplicationsImage struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key"`
	ApplicationID uuid.UUID `gorm:"type:char(36);not null"`
	MimeType      string    `gorm:"type:text;not null" json:"-"`
	CreatedAt     time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (image ApplicationsImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(image.ID)
}

type ApplicationsImageRepository interface {
	CreateApplicationsImage(applicationId uuid.UUID, src io.Reader, mimeType string) (ApplicationsImage, error)
	GetApplicationsImage(id uuid.UUID) (ApplicationsImage, error)
	OpenApplicationsImage(appImg ApplicationsImage) (io.ReadCloser, error)
	DeleteApplicationsImage(appImg ApplicationsImage) error
}

type applicationsImageRepository struct {
	storage storagePkg.Storage
}

func NewApplicationsImageRepository(storage storagePkg.Storage) ApplicationsImageRepository {
	return &applicationsImageRepository{storage: storage}
}

func (repo *applicationsImageRepository) CreateApplicationsImage(applicationId uuid.UUID, src io.Reader, mimeType string) (ApplicationsImage, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return ApplicationsImage{}, err
	}

	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return ApplicationsImage{}, err
	} else if len(ext) == 0 {
		return ApplicationsImage{}, fmt.Errorf("%s is not registered", mimeType)
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext[0])

	err = repo.storage.Save(filename, src)
	if err != nil {
		return ApplicationsImage{}, err
	}

	im := ApplicationsImage{
		ID:            id,
		ApplicationID: applicationId,
		MimeType:      mimeType,
	}

	if err = db.Create(&im).Error; err != nil {
		_ = repo.storage.Delete(filename)
		return ApplicationsImage{}, err
	}

	return im, nil
}

func (repo *applicationsImageRepository) GetApplicationsImage(id uuid.UUID) (ApplicationsImage, error) {
	im := ApplicationsImage{
		ID: id,
	}

	if err := db.First(&im).Error; err != nil {
		return ApplicationsImage{}, err
	}

	return im, nil
}

func (repo *applicationsImageRepository) OpenApplicationsImage(appImg ApplicationsImage) (io.ReadCloser, error) {
	ext, err := mime.ExtensionsByType(appImg.MimeType)
	if err != nil {
		return nil, err
	} else if len(ext) == 0 {
		return nil, fmt.Errorf("%s is not registered", appImg.MimeType)
	}

	filenames := []string{}
	for _, ex := range ext {
		filenames = append(filenames, fmt.Sprintf("%s%s", appImg.ID.String(), ex))
	}

	for _, filename := range filenames {
		rd, err := repo.storage.Open(filename)
		if err == nil {
			return rd, nil
		}
	}

	return nil, fmt.Errorf("image not found")
}

func (repo *applicationsImageRepository) DeleteApplicationsImage(appImg ApplicationsImage) error {
	ext, err := mime.ExtensionsByType(appImg.MimeType)
	if err != nil {
		return err
	} else if len(ext) == 0 {
		return fmt.Errorf("%s is not registered", appImg.MimeType)
	}

	filenames := []string{}
	for _, ex := range ext {
		filenames = append(filenames, fmt.Sprintf("%s%s", appImg.ID.String(), ex))
	}

	for _, filename := range filenames {
		err := repo.storage.Delete(filename)
		if err == nil {
			if err := db.Delete(appImg).Error; err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("image not found")
}
