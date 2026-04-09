package router

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/service"
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

func (h Handlers) PostFile(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(*service.User)
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error("failed to parse request as multipart/form-data", zap.Error(err))
		return service.NewUnexpectedError(err)
	}
	files, ok := form.File["file"]
	if !ok || len(files) != 1 {
		logger.Info("could not find field `file` in request, or its length is not 1")
		return service.NewBadInputError("invalid file")
	}
	reqfile := files[0]
	names, ok := form.Value["name"]
	if !ok || len(names) != 1 {
		logger.Info("could not find field `name` in request, or its length is not 1")
		return service.NewBadInputError("invalid file name")
	}
	name := names[0]
	applicationIDs, ok := form.Value["application_id"]
	if !ok || len(applicationIDs) != 1 {
		logger.Info("could not find field `application_id` in request, or its length is not 1")
		return service.NewBadInputError("invalid file application id")
	}
	applicationID, err := uuid.Parse(applicationIDs[0])
	if err != nil {
		logger.Info("could not parse application_id as UUID", zap.Error(err))
		return service.NewBadInputError("invalid file application id").
			WithInternal(err)
	}

	mimetype := reqfile.Header.Get(echo.HeaderContentType)

	src, err := reqfile.Open()
	if err != nil {
		logger.Error("failed to open requested file", zap.Error(err))
		return service.NewUnexpectedError(err)
	}
	defer src.Close()

	file, err := h.Service.WriteFile(ctx, loginUser, applicationID, name, mimetype, src)
	if err != nil {
		logger.Error("failed to write file in service", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusOK, &FileResponse{file.ID})
}

func (h Handlers) GetFile(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	fileID, err := uuid.Parse(c.Param("fileID"))
	if err != nil {
		logger.Error("could not parse query parameter `fileID` as UUID", zap.Error(err))
		return service.NewBadInputError("invalid file ID").
			WithInternal(err)
	}

	file, err := h.Service.GetFile(ctx, fileID)
	if err != nil {
		return err
	}

	modifiedAt := file.CreatedAt.Truncate(time.Second)

	im := c.Request().Header.Get(echo.HeaderIfModifiedSince)
	if im != "" {
		imt, err := http.ParseTime(im)
		if err != nil {
			logger.Info("could not parse time in request header", zap.Error(err))
			return service.NewBadInputError("invalid If-Modified-Since header").
				WithInternal(err)
		}
		if modifiedAt.Before(imt) || modifiedAt.Equal(imt) {
			logger.Info(
				"content is not modified since the last request",
				zap.String("ID", fileID.String()),
				zap.Time("If-Modified-Since", imt))
			return c.NoContent(http.StatusNotModified)
		}
	}

	content, err := h.Service.ReadFile(ctx, fileID)
	if err != nil {
		return err
	}
	defer content.Close()

	c.Response().Header().Set("Cache-Control", "private, no-cache, max-age=0")
	c.Response().Header().Set(echo.HeaderLastModified, modifiedAt.UTC().Format(http.TimeFormat))

	return c.Stream(http.StatusOK, file.MimeType, content)
}

func (h Handlers) GetFileMeta(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	fileID, err := uuid.Parse(c.Param("fileID"))
	if err != nil {
		logger.Info("could not parse query parameter `fileID` as UUID", zap.Error(err))
		return service.NewBadInputError("invalid file ID").
			WithInternal(err)
	}

	file, err := h.Service.GetFile(ctx, fileID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &FileMetaResponse{
		ID:        file.ID,
		Name:      file.Name,
		MimeType:  file.MimeType,
		CreatedBy: file.CreatedBy,
		CreatedAt: file.CreatedAt,
	})
}

func (h Handlers) DeleteFile(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(*service.User)
	fileID, err := uuid.Parse(c.Param("fileID"))
	if err != nil {
		logger.Info("could not parse query parameter `fileID` as UUID", zap.Error(err))
		return service.NewBadInputError("invalid file ID").
			WithInternal(err)
	}
	err = h.Service.DeleteFile(ctx, loginUser, fileID)
	if err != nil {
		return service.NewUnexpectedError(err)
	}

	return c.NoContent(http.StatusOK)
}
