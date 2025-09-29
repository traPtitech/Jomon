package router

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/internal/model"
	"github.com/traPtitech/Jomon/internal/testutil"
	"github.com/traPtitech/Jomon/internal/testutil/random"
	"go.uber.org/mock/gomock"
)

// nolint:lll
var testJpeg = `/9j/4AAQSkZJRgABAQIAOAA4AAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////2wBDAf//////////////////////////////////////////////////////////////////////////////////////wAARCAABAAEDAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwBKBH//2Q`

func TestHandlers_PostFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		requestID := uuid.New()
		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			defer func() {
				err := writer.Close()
				require.NoError(t, err)
			}()
			err := writer.WriteField("name", "test")
			require.NoError(t, err)
			err = writer.WriteField("request_id", requestID.String())
			require.NoError(t, err)
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", "image/jpeg")
			part.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
			pw, err := writer.CreatePart(part)
			require.NoError(t, err)
			file, err := base64.RawStdEncoding.DecodeString(testJpeg)
			require.NoError(t, err)
			_, err = pw.Write(file)
			require.NoError(t, err)
		}()

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files")
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), "test", "image/jpeg", requestID, user.ID).
			Return(file, nil)

		h.Storage.
			EXPECT().
			Save(c.Request().Context(), file.ID.String(), gomock.Any()).
			Return(nil)

		require.NoError(t, h.Handlers.PostFile(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedToRepositoryCreateFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		requestID := uuid.New()

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			defer func() {
				err := writer.Close()
				require.NoError(t, err)
			}()
			err := writer.WriteField("name", "test")
			require.NoError(t, err)
			err = writer.WriteField("request_id", requestID.String())
			require.NoError(t, err)
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", "image/jpeg")
			part.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
			pw, err := writer.CreatePart(part)
			require.NoError(t, err)
			file, err := base64.RawStdEncoding.DecodeString(testJpeg)
			require.NoError(t, err)
			_, err = pw.Write(file)
			require.NoError(t, err)
		}()

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files")
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		mocErr := errors.New("failed to create file")

		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), "test", "image/jpeg", requestID, user.ID).
			Return(nil, mocErr)

		err = h.Handlers.PostFile(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})

	t.Run("FailedToServiceCreateFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		requestID := uuid.New()
		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			defer func() {
				err := writer.Close()
				require.NoError(t, err)
			}()
			err := writer.WriteField("name", "test")
			require.NoError(t, err)
			err = writer.WriteField("request_id", requestID.String())
			require.NoError(t, err)
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", "image/jpeg")
			part.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
			pw, err := writer.CreatePart(part)
			require.NoError(t, err)
			file, err := base64.RawStdEncoding.DecodeString(testJpeg)
			require.NoError(t, err)
			_, err = pw.Write(file)
			require.NoError(t, err)
		}()

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/files", pr)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files")
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			CreateFile(c.Request().Context(), "test", "image/jpeg", requestID, user.ID).
			Return(file, nil)

		mocErr := errors.New("failed to save file")

		h.Storage.
			EXPECT().
			Save(c.Request().Context(), file.ID.String(), gomock.Any()).
			Return(mocErr)

		err = h.Handlers.PostFile(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})
}

func TestHandlers_GetFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		path := fmt.Sprintf("/api/files/%s", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)

		h.Storage.
			EXPECT().
			Open(c.Request().Context(), file.ID.String()).
			Return(r, nil)

		require.NoError(t, h.Handlers.GetFile(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedToGetFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		mocErr := errors.New("file not found")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(nil, mocErr)

		err = h.Handlers.GetFile(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})

	t.Run("FailedToOpenFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)

		mocErr := errors.New("failed to open file")

		h.Storage.
			EXPECT().
			Open(c.Request().Context(), file.ID.String()).
			Return(nil, mocErr)

		err = h.Handlers.GetFile(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/files/po", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues("po")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, mocErr := uuid.Parse("po")

		err = h.Handlers.GetFile(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, mocErr), err)
	})
}

func TestHandlers_GetFileMeta(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s/meta", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID/meta")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)

		require.NoError(t, h.Handlers.GetFileMeta(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var res *FileMetaResponse
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		exp := &FileMetaResponse{
			ID:        file.ID,
			Name:      file.Name,
			MimeType:  file.MimeType,
			CreatedBy: file.CreatedBy,
			CreatedAt: file.CreatedAt,
		}
		testutil.RequireEqual(t, exp, res)
	})

	t.Run("FailedToGetFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s/meta", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID/meta")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())

		mocErr := errors.New("file not found")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(nil, mocErr)

		err = h.Handlers.GetFileMeta(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/files/po/meta", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID/meta")
		c.SetParamNames("fileID")
		c.SetParamValues("po")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, mocErr := uuid.Parse("po")

		err = h.Handlers.GetFileMeta(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, mocErr), err)
	})
}

func TestHandlers_DeleteFile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)
		h.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(c.Request().Context(), file.ID).
			Return(nil)
		h.Storage.
			EXPECT().
			Delete(c.Request().Context(), file.ID.String()).
			Return(nil)

		require.NoError(t, h.Handlers.DeleteFile(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedToRepositoryDeleteFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())
		c.Set(loginUserKey, user)

		mocErr := errors.New("file could not be deleted")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)
		h.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(c.Request().Context(), file.ID).
			Return(mocErr)

		err = h.Handlers.DeleteFile(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})

	t.Run("FailedToServiceDeleteFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		file := &model.File{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			MimeType:  "image/jpeg",
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s", file.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(file.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockFileRepository.
			EXPECT().
			GetFile(c.Request().Context(), file.ID).
			Return(file, nil)
		h.Repository.MockFileRepository.
			EXPECT().
			DeleteFile(c.Request().Context(), file.ID).
			Return(nil)

		mocErr := errors.New("failed to delete file")

		h.Storage.
			EXPECT().
			Delete(c.Request().Context(), file.ID.String()).
			Return(mocErr)

		err = h.Handlers.DeleteFile(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		invalidUUID := "invalid-uuid"
		_, mocErr := uuid.Parse(invalidUUID)

		e := echo.New()
		path := fmt.Sprintf("/api/files/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/files/:fileID")
		c.SetParamNames("fileID")
		c.SetParamValues(invalidUUID)
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteFile(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, mocErr), err)
	})
}
