package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

const (
	contextUserKey        = "user"
	contextAccessTokenKey = "access_token"
)

func (s *Service) GetUsers(c echo.Context) error {
	token := c.Get(contextAccessTokenKey).(string)

	users, err := s.Users.GetUsers(token)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	for i := range users {
		users[i].GiveIsUserAdmin(admins)
	}

	return c.JSON(http.StatusOK, users)
}

func (s *Service) GetMyUser(c echo.Context) error {
	user, ok := c.Get(contextUserKey).(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Service) GetAdminUsers(c echo.Context) error {
	token := c.Get(contextAccessTokenKey).(string)

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

	myUser, ok := c.Get(contextUserKey).(model.User)
	if !ok || myUser.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	if !myUser.IsAdmin {
		return c.NoContent(http.StatusForbidden)
	}

	token := c.Get(contextAccessTokenKey).(string)
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
