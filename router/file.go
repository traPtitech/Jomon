package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FileResponse struct {
	FileIDs []*uuid.UUID `json:"file_ids"`
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

func (h *Handlers) PostFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return internalServerError(err)
	}
	var fileIDs []*uuid.UUID
	files, ok := form.File["file"]
	if !ok || len(files) != 1 {
		return badRequest(fmt.Errorf("invalid file"))
	}
	file := files[0]
	names, ok := form.Value["name"]
	if !ok || len(names) != 1 {
		return badRequest(fmt.Errorf("invalid file name"))
	}
	name := names[0]
	requestIDs, ok := form.Value["request_id"]
	if !ok || len(requestIDs) != 1 {
		return badRequest(fmt.Errorf("invalid file request id"))
	}
	requestID, err := uuid.Parse(requestIDs[0])
	if err != nil {
		return badRequest(err)
	}

	mimetype := file.Header.Get(echo.HeaderContentType)
	if !acceptedMimeTypes[mimetype] {
		return c.NoContent(http.StatusUnsupportedMediaType)
	}

	src, err := file.Open()
	if err != nil {
		return internalServerError(err)
	}
	defer src.Close()

	fileID, err := h.Service.CreateFile(src, name, mimetype, requestID)
	if err != nil {
		return internalServerError(err)
	}
	fileIDs = append(fileIDs, &fileID.ID)
	return c.JSON(http.StatusOK, &FileResponse{fileIDs})
}

func (h *Handlers) GetFile(c echo.Context) error {
	fileID, err := uuid.Parse(c.Param("fileID"))
	if err != nil {
		return badRequest(err)
	}
	file, err := h.Service.GetFile(fileID)
	if err != nil {
		internalServerError(err)
	}

	modifiedAt := file.CreatedAt.Truncate(time.Second)

	im := c.Request().Header.Get(echo.HeaderIfModifiedSince)
	if im != "" {
		imt, err := http.ParseTime(im)
		if err != nil {
			return badRequest(err)
		}
		if modifiedAt.Before(imt) || modifiedAt.Equal(imt) {
			return c.NoContent(http.StatusNotModified)
		}
	}
	f, err := h.Service.OpenFile(fileID)
	if err != nil {
		return internalServerError(err)
	}
	defer f.Close()

	c.Response().Header().Set("Cache-Control", "private, no-cache, max-age=0")
	c.Response().Header().Set(echo.HeaderLastModified, modifiedAt.UTC().Format(http.TimeFormat))

	return c.Stream(http.StatusOK, file.MimeType, f)
}

func (h *Handlers) DeleteFile(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
