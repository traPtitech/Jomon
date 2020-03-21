package router

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

func (s *Service) GetUsers(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")

	users, err := s.Users.GetUsers(token)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	for _, user := range users {
		user.GiveIsUserAdmin(admins)
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

	allUsers, err := s.Users.GetUsers(token)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	adminUsers := []model.User{}
	for _, user := range allUsers {
		user.GiveIsUserAdmin(admins)
		if user.IsAdmin {
			adminUsers = append(adminUsers, user)
		}
	}

	return c.JSON(http.StatusOK, adminUsers)
}

type PutAdminRequest struct {
	TrapId  string `json:"trap_id"`
	ToAdmin bool   `json:"to_admin"`
}

func (s *Service) PutAdminUsers(c echo.Context) error {
	var req PutAdminRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	myUser, ok := c.Get("user").(model.User)
	if !ok || myUser.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	if !myUser.IsAdmin {
		return c.NoContent(http.StatusForbidden)
	}

	token := c.Request().Header.Get("Authorization")
	exist, err := s.Users.ExistsUser(token, req.TrapId)
	if !exist || err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user := model.User{
		TrapId: req.TrapId,
	}
	user.GiveIsUserAdmin(admins)

	if req.ToAdmin == user.IsAdmin {
		return c.NoContent(http.StatusConflict)
	}

	if req.ToAdmin {
		if err := s.Administrators.AddAdministrator(req.TrapId); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	} else {
		if err := s.Administrators.RemoveAdministrator(req.TrapId); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	admins, err = s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	user.GiveIsUserAdmin(admins)

	return c.JSON(http.StatusOK, user)
}

func (s *Service) SetMyUser(c echo.Context) (echo.Context, error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c, errors.New("no token")
	}

	user, err := s.Users.GetMyUser(token)
	if err != nil {
		return c, err
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c, err
	}

	user.GiveIsUserAdmin(admins)

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
