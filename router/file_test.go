package router

import (
	"bytes"
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
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil/random"
)

var testJpeg = `/9j/4AAQSkZJRgABAQIAOAA4AAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////2wBDAf//////////////////////////////////////////////////////////////////////////////////////wAARCAABAAEDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwBKBH//2Q`

func TestHandlers_PostFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		request := uuid.New()
		user := User{
			ID:          uuid.New(),
			Name:        "test",
			DisplayName: "test",
			Admin:       true,
		}

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
		req, err := http.NewRequest(http.MethodPost, "/api/files", pr)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		require.NoError(t, err)
		sess.Values[sessionUserKey] = user
		require.NoError(t, sess.Save(c.Request(), c.Response()))

		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), "test", "image/jpeg", request, user.ID).
			Return(file, nil)

		h.Storage.
			EXPECT().
			Save(file.ID.String(), gomock.Any()).
			Return(nil)

		if assert.NoError(t, h.Handlers.PostFile(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedToRepositoryCreateFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		request := uuid.New()
		user := User{
			ID:          uuid.New(),
			Name:        "test",
			DisplayName: "test",
			Admin:       true,
		}

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
		req, err := http.NewRequest(http.MethodPost, "/api/files", pr)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		require.NoError(t, err)
		sess.Values[sessionUserKey] = user
		require.NoError(t, sess.Save(c.Request(), c.Response()))

		mocErr := errors.New("failed to create file")

		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), "test", "image/jpeg", request, user.ID).
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
		user := User{
			ID:          uuid.New(),
			Name:        "test",
			DisplayName: "test",
			Admin:       true,
		}

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
		req, err := http.NewRequest(http.MethodPost, "/api/files", pr)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		require.NoError(t, err)
		sess.Values[sessionUserKey] = user
		require.NoError(t, sess.Save(c.Request(), c.Response()))

		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), "test", "image/jpeg", request, user.ID).
			Return(file, nil)

		mocErr := errors.New("failed to save file")

		h.Storage.
			EXPECT().
			Save(file.ID.String(), gomock.Any()).
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
			Name:      random.AlphaNumeric(t, 20),
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

		h.Storage.
			EXPECT().
			Open(file.ID.String()).
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
			Name:      random.AlphaNumeric(t, 20),
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
			Name:      random.AlphaNumeric(t, 20),
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

		mocErr := errors.New("failed to open file")

		h.Storage.
			EXPECT().
			Open(file.ID.String()).
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

func TestHandlers_GetFileMeta(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/files/%s/meta", file.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID/meta")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)

		if assert.NoError(t, h.Handlers.GetFileMeta(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedToGetFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/files/%s/meta", file.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID/meta")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		mocErr := errors.New("file not found")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(nil, mocErr)

		err = h.Handlers.GetFileMeta(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/files/po/meta", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID/meta")
		c.SetParamNames("fileID")
		c.SetParamValues("po")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		_, mocErr := uuid.Parse("po")

		err = h.Handlers.GetFileMeta(c)
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
			Name:      random.AlphaNumeric(t, 20),
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
			Return(nil)

		h.Storage.
			EXPECT().
			Delete(file.ID.String()).
			Return(nil)

		if assert.NoError(t, h.Handlers.DeleteFile(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedToRepositoryDeleteFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
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

		mocErr := errors.New("file could not be deleted")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(c.Request().Context(), file.ID).
			Return(mocErr)

		err = h.Handlers.DeleteFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("FailedToServiceDeleteFile", func(t *testing.T) {
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
			Return(nil)

		mocErr := errors.New("failed to delete file")

		h.Storage.
			EXPECT().
			Delete(file.ID.String()).
			Return(mocErr)

		err = h.Handlers.DeleteFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodDelete, "/api/files/po", nil)
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

		err = h.Handlers.DeleteFile(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, mocErr), err)
		}
	})
}
