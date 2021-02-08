package model

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type storageMock struct {
	mock.Mock
}

func (m *storageMock) Save(filename string, src io.Reader) error {
	ret := m.Called(filename, src)
	return ret.Error(0)
}

func (m *storageMock) Open(filename string) (io.ReadCloser, error) {
	ret := m.Called(filename)
	return ret.Get(0).(io.ReadCloser), ret.Error(1)
}

func (m *storageMock) Delete(filename string) error {
	ret := m.Called(filename)
	return ret.Error(0)
}

func TestCreateApplicationsImage(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		var actualFilename string
		var actualReaderString string

		sm := new(storageMock)
		sm.On("Save", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			actualFilename = args.String(0)
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(args.Get(1).(io.Reader))
			actualReaderString = buf.String()
		}).Return(nil).Once()

		imageRepo := NewFileRepository(sm)

		sampleText := "sampleData"

		appId, err := repo.createRequest(db, "")
		if err != nil {
			panic(err)
		}

		mimeType := "image/png"

		im, err := imageRepo.CreateFile(appId, strings.NewReader(sampleText), mimeType)
		asr.NoError(err)
		asr.Equal(fmt.Sprintf("%s.png", im.ID.String()), actualFilename)
		asr.Equal(sampleText, actualReaderString)
		asr.Equal(mimeType, im.MimeType)
		asr.NotZero(im.CreatedAt)
	})
}

func TestGetApplicationsImage(t *testing.T) {
	t.Parallel()

	sm := new(storageMock)
	sm.On("Save", mock.Anything, mock.Anything).Return(nil)

	imageRepo := NewFileRepository(sm)

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createRequest(db, "")
		if err != nil {
			panic(err)
		}

		createdIm, err := imageRepo.CreateFile(appId, strings.NewReader(""), "image/png")
		if err != nil {
			panic(err)
		}

		getIm, err := imageRepo.GetFile(createdIm.ID)
		asr.NoError(err)
		asr.Equal(createdIm.ID, getIm.ID)
		asr.Equal(createdIm.RequestID, getIm.RequestID)
		asr.Equal(createdIm.MimeType, getIm.MimeType)
		asr.WithinDuration(createdIm.CreatedAt, getIm.CreatedAt, 1*time.Second)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = imageRepo.GetFile(id)
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})
}

func TestDeleteApplicationsImage(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		sm := new(storageMock)
		sm.On("Save", mock.Anything, mock.Anything).Return(nil)
		sm.On("Delete", mock.Anything).Return(nil).Once()

		imageRepo := NewFileRepository(sm)

		appId, err := repo.createRequest(db, "")
		if err != nil {
			panic(err)
		}

		im, err := imageRepo.CreateFile(appId, strings.NewReader(""), "image/png")
		if err != nil {
			panic(err)
		}

		err = imageRepo.DeleteFile(im)
		asr.NoError(err)

		_, err = imageRepo.GetFile(im.ID)
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})
}
