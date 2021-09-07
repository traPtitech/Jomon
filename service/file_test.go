package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testJpeg = `/9j/4AAQSkZJRgABAQIAOAA4AAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////2wBDAf//////////////////////////////////////////////////////////////////////////////////////wAARCAABAAEDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwBKBH//2Q`

func TestService_CreateFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file, err := base64.RawStdEncoding.DecodeString(testJpeg)
		require.NoError(t, err)
		r := bytes.NewReader(file)

		fileID := uuid.New()
		mimetype := "image/jpeg"

		ext, err := mime.ExtensionsByType(mimetype)
		require.NoError(t, err)

		strg := NewMockStorage(ctrl)
		strg.EXPECT().
			Save(fmt.Sprintf("%s%s", fileID.String(), ext[0]), r).
			Return(nil)
		s, err := NewServices(strg)
		require.NoError(t, err)

		err = s.CreateFile(r, fileID, mimetype)
		assert.NoError(t, err)
	})

	t.Run("FailedToSave", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file, err := base64.RawStdEncoding.DecodeString(testJpeg)
		require.NoError(t, err)
		r := bytes.NewReader(file)

		fileID := uuid.New()
		mimetype := "image/jpeg"

		ext, err := mime.ExtensionsByType(mimetype)
		require.NoError(t, err)

		strg := NewMockStorage(ctrl)
		strg.EXPECT().
			Save(fmt.Sprintf("%s%s", fileID.String(), ext[0]), r).
			Return(errors.New("failed to save file"))
		s, err := NewServices(strg)
		require.NoError(t, err)

		err = s.CreateFile(r, fileID, mimetype)
		assert.Error(t, err)
	})

	t.Run("UnknownMimetype", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file, err := base64.RawStdEncoding.DecodeString(testJpeg)
		require.NoError(t, err)
		r := bytes.NewReader(file)

		fileID := uuid.New()
		mimetype := "po"

		strg := NewMockStorage(ctrl)
		s, err := NewServices(strg)
		require.NoError(t, err)

		err = s.CreateFile(r, fileID, mimetype)
		assert.Error(t, err)
	})
}

func TestService_OpenFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file, err := base64.RawStdEncoding.DecodeString(testJpeg)
		require.NoError(t, err)
		r := bytes.NewReader(file)
		rc := io.NopCloser(r)

		fileID := uuid.New()
		mimetype := "image/jpeg"

		ext, err := mime.ExtensionsByType(mimetype)
		require.NoError(t, err)

		strg := NewMockStorage(ctrl)
		strg.EXPECT().
			Open(fmt.Sprintf("%s%s", fileID.String(), ext[0])).
			Return(rc, nil)
		s, err := NewServices(strg)
		require.NoError(t, err)

		got, err := s.OpenFile(fileID, mimetype)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("FailedToOpen", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		fileID := uuid.New()
		mimetype := "image/jpeg"

		ext, err := mime.ExtensionsByType(mimetype)
		require.NoError(t, err)

		strg := NewMockStorage(ctrl)
		strg.EXPECT().
			Open(fmt.Sprintf("%s%s", fileID.String(), ext[0])).
			Return(nil, errors.New("failed to open file"))
		s, err := NewServices(strg)
		require.NoError(t, err)

		_, err = s.OpenFile(fileID, mimetype)
		assert.Error(t, err)
	})

	t.Run("UnknownMimetype", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		fileID := uuid.New()
		mimetype := "po"

		strg := NewMockStorage(ctrl)
		s, err := NewServices(strg)
		require.NoError(t, err)

		_, err = s.OpenFile(fileID, mimetype)
		assert.Error(t, err)
	})
}

func TestService_DeleteFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		fileID := uuid.New()
		mimetype := "image/jpeg"

		ext, err := mime.ExtensionsByType(mimetype)
		require.NoError(t, err)

		strg := NewMockStorage(ctrl)
		strg.EXPECT().
			Delete(fmt.Sprintf("%s%s", fileID.String(), ext[0])).
			Return(nil)
		s, err := NewServices(strg)
		require.NoError(t, err)

		err = s.DeleteFile(fileID, mimetype)
		assert.NoError(t, err)
	})

	t.Run("FailedToDelete", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		fileID := uuid.New()
		mimetype := "image/jpeg"

		ext, err := mime.ExtensionsByType(mimetype)
		require.NoError(t, err)

		strg := NewMockStorage(ctrl)
		strg.EXPECT().
			Delete(fmt.Sprintf("%s%s", fileID.String(), ext[0])).
			Return(errors.New("failed to delete file"))
		s, err := NewServices(strg)
		require.NoError(t, err)

		err = s.DeleteFile(fileID, mimetype)
		assert.Error(t, err)
	})

	t.Run("UnknownMimetype", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		fileID := uuid.New()
		mimetype := "po"

		strg := NewMockStorage(ctrl)
		s, err := NewServices(strg)
		require.NoError(t, err)

		err = s.DeleteFile(fileID, mimetype)
		assert.Error(t, err)
	})
}