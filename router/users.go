package router

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

func (s *Service) GetUsers(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	users, err := s.Users.GetUsers(token, admins, false)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, users)
}

func (s *Service) GetMyUser(c echo.Context) error {
	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Service) GetAdminUsers(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	users, err := s.Users.GetUsers(token, admins, true)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, users)
}

type PutAdminRequest struct {
	trapId  string `json:"trap_id"`
	toAdmin bool   `json:"to_admin"`
}

func (s *Service) PutAdminUsers(c echo.Context) error {
	var req PutAdminRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	//403
	//404

	if req.toAdmin {
		s.Administrators.AddAdministrator(req.trapId)
	} else {
		s.Administrators.RemoveAdministrator(req.trapId)
	}

	return c.JSON(http.StatusOK, req)
}

func (s *Service) SetMyUser(c echo.Context) (echo.Context, error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c, errors.New("no token")
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c, err
	}

	user, err := s.Users.GetMyUser(token, admins)
	if err != nil {
		return c, err
	}

	c.Set("user", user)

	return c, nil
}

func (s *Service) AuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c, err := s.SetMyUser(c)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		return next(c)
	}
}
