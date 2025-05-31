package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/logging"
	"go.uber.org/zap"
)

type FileResponse struct {
	ID uuid.UUID `json:"id"`
}

type FileMetaResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	MimeType  string    `json:"mime_type"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

var acceptedMimeTypes = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"image/bmp":          true,
	"application/pdf":    true,
	"application/msword": true,
	"application/zip":    true,
}

func (h Handlers) PostFile(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	form, err := c.MultipartForm()
	if err != nil {
		logger.Error("failed to parse request as multipart/form-data", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	files, ok := form.File["file"]
	if !ok || len(files) != 1 {
		logger.Info("could not find field `file` in request, or its length is not 1")
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid file"))
	}
	reqfile := files[0]
	names, ok := form.Value["name"]
	if !ok || len(names) != 1 {
		logger.Info("could not find field `name` in request, or its length is not 1")
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid file name"))
	}
	name := names[0]
	requestIDs, ok := form.Value["request_id"]
	if !ok || len(requestIDs) != 1 {
		logger.Info("could not find field `request_id` in request, or its length is not 1")
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid file request id"))
	}
	requestID, err := uuid.Parse(requestIDs[0])
	if err != nil {
		logger.Info("could not parse request_id as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	mimetype := reqfile.Header.Get(echo.HeaderContentType)
	if !acceptedMimeTypes[mimetype] {
		logger.Info("requested unsupported mime type", zap.String("mime-type", mimetype))
		return echo.NewHTTPError(
			http.StatusUnsupportedMediaType,
			fmt.Errorf("unsupported media type"))
	}

	src, err := reqfile.Open()
	if err != nil {
		logger.Error("failed to open requested file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	// get create user
	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		logger.Error("failed to get session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user, ok := sess.Values[sessionUserKey].(User)
	if !ok {
		logger.Error("failed to parse stored session as user info")
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid user"))
	}

	file, err := h.Repository.CreateFile(ctx, name, mimetype, requestID, user.ID)
	if err != nil {
		logger.Error("failed to create file in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = h.Storage.Save(file.ID.String(), src)
	if err != nil {
		logger.Error("failed to save file id in storage", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &FileResponse{file.ID})
}

func (h Handlers) GetFile(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	fileID, err := uuid.Parse(c.Param("fileID"))
	if err != nil {
		logger.Error("could not parse query parameter `fileID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	file, err := h.Repository.GetFile(ctx, fileID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info("could not find file in repository", zap.String("ID", fileID.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to get file from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	modifiedAt := file.CreatedAt.Truncate(time.Second)

	im := c.Request().Header.Get(echo.HeaderIfModifiedSince)
	if im != "" {
		imt, err := http.ParseTime(im)
		if err != nil {
			logger.Info("could not parse time in request header", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if modifiedAt.Before(imt) || modifiedAt.Equal(imt) {
			logger.Info(
				"content is not modified since the last request",
				zap.String("ID", fileID.String()),
				zap.Time("If-Modified-Since", imt))
			return c.NoContent(http.StatusNotModified)
		}
	}

	f, err := h.Storage.Open(fileID.String())
	if err != nil {
		logger.Error(
			"failed to open file in storage",
			zap.String("ID", fileID.String()),
			zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer f.Close()

	c.Response().Header().Set("Cache-Control", "private, no-cache, max-age=0")
	c.Response().Header().Set(echo.HeaderLastModified, modifiedAt.UTC().Format(http.TimeFormat))

	return c.Stream(http.StatusOK, file.MimeType, f)
}

func (h Handlers) GetFileMeta(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	fileID, err := uuid.Parse(c.Param("fileID"))
	if err != nil {
		logger.Info("could not parse query parameter `fileID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	file, err := h.Repository.GetFile(ctx, fileID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info("could not find file in repository", zap.String("ID", fileID.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to get file from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &FileMetaResponse{
		ID:        file.ID,
		Name:      file.Name,
		MimeType:  file.MimeType,
		CreatedBy: file.CreatedBy,
		CreatedAt: file.CreatedAt,
	})
}

func (h Handlers) DeleteFile(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	fileID, err := uuid.Parse(c.Param("fileID"))
	if err != nil {
		logger.Info("could not parse query parameter `fileID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.Repository.DeleteFile(ctx, fileID)
	if err != nil {
		if ent.IsConstraintError(err) {
			logger.Info("constraint error while deleting file", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		logger.Error("failed to delete file in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = h.Storage.Delete(fileID.String())
	if err != nil {
		logger.Error("failed to delete file in storage", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// isFileCreator 与えられたユーザーがファイルの作成者かどうかを確認します
func (h Handlers) isFileCreator(ctx context.Context, userID, fileID uuid.UUID) (bool, error) {
	file, err := h.Repository.GetFile(ctx, fileID)
	if err != nil {
		return false, err
	}
	return file.CreatedBy == userID, nil
}

func (h Handlers) filterAdminOrFileCreator(
	ctx context.Context, userID uuid.UUID, fileID uuid.UUID,
) *echo.HTTPError {
	logger := logging.GetLogger(ctx)
	isAdmin, err := h.isAdmin(ctx, userID)
	if err != nil {
		// NOTE: end.IsNotFound(err)は起こり得ないと仮定
		logger.Error("failed to check if user is admin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if isAdmin {
		return nil
	}
	isCreator, err := h.isFileCreator(ctx, userID, fileID)
	if err != nil {
		logger.Error("failed to check if user is file creator", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if isCreator {
		return nil
	}
	return echo.NewHTTPError(http.StatusForbidden, errUserIsNotAdminOrFileCreator)
}
