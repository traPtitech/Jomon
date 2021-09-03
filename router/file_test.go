package router

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
)

var testJpeg = `/9j/4AAQSkZJRgABAQIAOAA4AAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////2wBDAf//////////////////////////////////////////////////////////////////////////////////////wAARCAABAAEDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwBKBH//2Q`

func TestHandlers_PostFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		request := uuid.New()

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			defer writer.Close()
			writer.WriteField("name", "test")
			writer.WriteField("request_id", request.String())
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", "image/jpeg")
			part.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
			pw, err := writer.CreatePart(part)
			assert.NoError(t, err)
			file, err := base64.RawStdEncoding.DecodeString(testJpeg)
			assert.NoError(t, err)
			_, err = pw.Write(file)
			assert.NoError(t, err)
		}()

		file := &model.File{
			ID: uuid.New(),
		}

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			CreateFile(ctx, gomock.Any(), "test", "image/jpeg", request).
			Return(file, nil)
		th.Service.MockService.
			EXPECT().
			CreateFile(gomock.Any(), file.ID, "image/jpeg").
			Return(nil)

		req := httptest.NewRequest(echo.POST, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		th.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedToRepositoryCreateFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		request := uuid.New()

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			defer writer.Close()
			writer.WriteField("name", "test")
			writer.WriteField("request_id", request.String())
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", "image/jpeg")
			part.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
			pw, err := writer.CreatePart(part)
			assert.NoError(t, err)
			file, err := base64.RawStdEncoding.DecodeString(testJpeg)
			assert.NoError(t, err)
			_, err = pw.Write(file)
			assert.NoError(t, err)
		}()

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			CreateFile(ctx, gomock.Any(), "test", "image/jpeg", request).
			Return(nil, errors.New("failed to create file"))

		req := httptest.NewRequest(echo.POST, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		th.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("FailedToServiceCreateFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		request := uuid.New()

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			defer writer.Close()
			writer.WriteField("name", "test")
			writer.WriteField("request_id", request.String())
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", "image/jpeg")
			part.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
			pw, err := writer.CreatePart(part)
			assert.NoError(t, err)
			file, err := base64.RawStdEncoding.DecodeString(testJpeg)
			assert.NoError(t, err)
			_, err = pw.Write(file)
			assert.NoError(t, err)
		}()

		file := &model.File{
			ID: uuid.New(),
		}

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			CreateFile(ctx, gomock.Any(), "test", "image/jpeg", request).
			Return(file, nil)
		th.Service.MockService.
			EXPECT().
			CreateFile(gomock.Any(), file.ID, "image/jpeg").
			Return(errors.New("failed to create file"))

		req := httptest.NewRequest(echo.POST, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		th.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandlers_GetFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		f, err := base64.RawStdEncoding.DecodeString(testJpeg)
		require.NoError(t, err)
		r := io.NopCloser(bytes.NewReader(f))

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			GetFile(ctx, file.ID).
			Return(file, nil)
		th.Service.MockService.
			EXPECT().
			OpenFile(file.ID, file.MimeType).
			Return(r, nil)

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, fmt.Sprintf("/api/files/%s", file.ID), nil, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("FailedToGetFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			GetFile(ctx, file.ID).
			Return(nil, errors.New("file not found"))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, fmt.Sprintf("/api/files/%s", file.ID), nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("FailedToOpenFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			GetFile(ctx, file.ID).
			Return(file, nil)
		th.Service.MockService.
			EXPECT().
			OpenFile(file.ID, file.MimeType).
			Return(nil, errors.New("file could not be opened"))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, fmt.Sprintf("/api/files/%s", file.ID), nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

func TestHandlers_DeleteFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(ctx, file.ID).
			Return(file, nil)

		th.Service.MockService.
			EXPECT().
			DeleteFile(file.ID, file.MimeType).
			Return(nil)

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, fmt.Sprintf("/api/files/%s", file.ID), nil, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("FailedToRepositoryDeleteFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(ctx, file.ID).
			Return(nil, errors.New("file could not be deleted"))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, fmt.Sprintf("/api/files/%s", file.ID), nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("FailedToServiceDeleteFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		ctx := context.Background()
		th.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(ctx, file.ID).
			Return(file, nil)

		th.Service.MockService.
			EXPECT().
			DeleteFile(file.ID, file.MimeType).
			Return(errors.New("file could not be deleted"))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, fmt.Sprintf("/api/files/%s", file.ID), nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
