package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FileResponse struct {
	FileIDs []*uuid.UUID `json:"file_ids"`
}

var acceptedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/bmp":  true,
}

func (h *Handlers) PostFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var fileIDs []*uuid.UUID
	file := form.File["file"][0]
	name := form.Value["name"][0]
	requestID, err := uuid.Parse(form.Value["request_id"][0])
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	mimetype := file.Header.Get(echo.HeaderContentType)
	if !acceptedMimeTypes[mimetype] {
		return c.NoContent(http.StatusUnsupportedMediaType)
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer src.Close()

	fileID, err := h.Service.CreateFile(src, name, mimetype, requestID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	fileIDs = append(fileIDs, &fileID.ID)
	return c.JSON(http.StatusOK, &FileResponse{fileIDs})
}

func (h *Handlers) GetFile(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) DeleteFile(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
