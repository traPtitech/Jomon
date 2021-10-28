package router

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
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

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), gomock.Any(), "test", "image/jpeg", request).
			Return(file, nil)

		ext, err := mime.ExtensionsByType("image/jpeg")
		require.NoError(t, err)

		filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

		h.Storage.
			EXPECT().
			Save(filename, gomock.Any()).
			Return(nil)

		if assert.NoError(t, h.Handlers.PostFile(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedToRepositoryCreateFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

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

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mocErr := errors.New("failed to create file")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), gomock.Any(), "test", "image/jpeg", request).
			Return(nil, mocErr)

		err = h.Handlers.PostFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("FailedToServiceCreateFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

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

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), gomock.Any(), "test", "image/jpeg", request).
			Return(file, nil)

		ext, err := mime.ExtensionsByType("image/jpeg")
		require.NoError(t, err)

		filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

		mocErr := errors.New("failed to save file")

		h.Storage.
			EXPECT().
			Save(filename, gomock.Any()).
			Return(mocErr)

		err = h.Handlers.PostFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})
}

func TestHandlers_GetFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		f, err := base64.RawStdEncoding.DecodeString(testJpeg)
		require.NoError(t, err)
		r := io.NopCloser(bytes.NewReader(f))

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/files/%s", file.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)

		ext, err := mime.ExtensionsByType(file.MimeType)
		require.NoError(t, err)

		filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

		h.Storage.
			EXPECT().
			Open(filename).
			Return(r, nil)

		if assert.NoError(t, h.Handlers.GetFile(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedToGetFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/files/%s", file.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		mocErr := errors.New("file not found")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(nil, mocErr)

		err = h.Handlers.GetFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("FailedToOpenFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/files/%s", file.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)

		ext, err := mime.ExtensionsByType(file.MimeType)
		require.NoError(t, err)

		filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

		mocErr := errors.New("failed to open file")

		h.Storage.
			EXPECT().
			Open(filename).
			Return(nil, mocErr)

		err = h.Handlers.GetFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/files/po", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues("po")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		_, mocErr := uuid.Parse("po")

		err = h.Handlers.GetFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, mocErr), err)
		}
	})
}

func TestHandlers_DeleteFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/files/%s", file.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(c.Request().Context(), file.ID).
			Return(file, nil)

		ext, err := mime.ExtensionsByType(file.MimeType)
		require.NoError(t, err)

		filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

		h.Storage.
			EXPECT().
			Delete(filename).
			Return(nil)

		if assert.NoError(t, h.Handlers.DeleteFile(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	/*
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

			mimetype := file.MimeType
			ext, err := mime.ExtensionsByType(mimetype)
			require.NoError(t, err)

			filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

			th.Storage.MockStorage.
				EXPECT().
				Delete(filename).
				Return(errors.New("failed to delete file"))

			statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, fmt.Sprintf("/api/files/%s", file.ID), nil, nil)
			assert.Equal(t, http.StatusInternalServerError, statusCode)
		})

		t.Run("UnknownFile", func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			accessUser := mustMakeUser(t, false)
			th, err := SetupTestHandlers(t, ctrl, accessUser)
			assert.NoError(t, err)

			statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, "/api/files/po", nil, nil)
			assert.Equal(t, http.StatusBadRequest, statusCode)
		})
	*/
}
