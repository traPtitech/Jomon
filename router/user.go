package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Admin       bool      `json:"admin"`
}

func (h *Handlers) GetUsers(c echo.Context) error {
	users, err := h.Repository.GetUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := make([]User, 0, len(users))
	for _, user := range users {
		res = append(res, User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Admin:       user.Admin,
		})
	}

	return c.JSON(http.StatusOK, res)
}

type PutUserRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Admin       bool   `json:"admin"`
}

func (h *Handlers) UpdateUserInfo(c echo.Context) error {
	var newUser PutUserRequest
	if err := c.Bind(&newUser); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := h.Repository.GetUserByName(c.Request().Context(), newUser.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updated, err := h.Repository.UpdateUser(c.Request().Context(), user.ID, newUser.Name, newUser.DisplayName, newUser.Admin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, User{
		ID:          user.ID,
		Name:        updated.Name,
		DisplayName: updated.DisplayName,
		Admin:       updated.Admin,
	})
}

func (h *Handlers) GetMe(c echo.Context) error {
	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user, ok := sess.Values[sessionUserKey].(*User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user info")
	}

	return c.JSON(http.StatusOK, user)
}
