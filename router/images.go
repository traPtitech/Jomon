package router

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/traPtitech/Jomon/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Service) GetImages(c echo.Context) error {
	imageId, err := uuid.FromString(c.Param("imageId"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	image, err := s.Images.GetApplicationsImage(imageId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c.NoContent(http.StatusNotFound)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	modifiedAt := image.CreatedAt.Truncate(time.Second)

	im := c.Request().Header.Get("If-Modified-Since")
	if im != "" {
		imt, err := http.ParseTime(im)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		if modifiedAt.Before(imt) || modifiedAt.Equal(imt) {
			return c.NoContent(http.StatusNotModified)
		}
	}

	f, err := s.Images.OpenApplicationsImage(image)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	defer f.Close()

	c.Response().Header().Set("Cache-Control", "private, no-cache, max-age=0")
	c.Response().Header().Set("Last-Modified", modifiedAt.UTC().Format(http.TimeFormat))

	return c.Stream(http.StatusOK, image.MimeType, f)
}

func (s *Service) DeleteImages(c echo.Context) error {
	imageId, err := uuid.FromString(c.Param("imageId"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	image, err := s.Images.GetApplicationsImage(imageId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c.NoContent(http.StatusNotFound)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	app, err := s.Applications.GetApplication(image.ApplicationID, false)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	if user != app.CreateUserTrapID && !user.IsAdmin {
		return c.NoContent(http.StatusForbidden)
	}

	if err = s.Images.DeleteApplicationsImage(image); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
