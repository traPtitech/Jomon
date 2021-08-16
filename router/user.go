package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/service"
)

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Admin       bool   `json:"admin"`
}

type UserOverview struct {
	Name        string     `json:"name"`
	DisplayName string     `json:"display_name"`
	Admin       bool       `json:"admin"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type UserResponse struct {
	Users []*UserOverview `json:"users"`
}

func (h *Handlers) GetUsers(c echo.Context) error {
	ctx := context.Background()
	users, err := h.Repository.GetUsers(ctx)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	res := []*UserOverview{}
	for _, user := range users {
		res = append(res, &UserOverview{
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Admin:       user.Admin,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			DeletedAt:   user.DeletedAt,
		})
	}
	return c.JSON(http.StatusOK, &UserResponse{res})
}

func (h *Handlers) PutUser(c echo.Context) error {
	var updateUser UserOverview
	if err := c.Bind(&updateUser); err != nil {
		return badRequest(err)
	}
	ctx := context.Background()
	user, err := h.Repository.GetUserByName(ctx, updateUser.Name)
	if err != nil {
		return internalServerError(err)
	}
	res, err := h.Repository.UpdateUser(ctx, user.ID, updateUser.Name, updateUser.DisplayName, updateUser.Admin)
	if err != nil {
		return internalServerError(err)
	}
	return c.JSON(http.StatusOK, &UserOverview{
		Name:        res.Name,
		DisplayName: res.DisplayName,
		Admin:       res.Admin,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
		DeletedAt:   res.DeletedAt,
	})
}

func (h *Handlers) GetMe(c echo.Context) error {
	sess, err := h.SessionStore.Get(c.Request(), h.SessionName)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	traqUser, ok := sess.Values[sessionUserKey].(*service.User)
	if !ok {
		c.Logger().Error(errors.New("failed to get users."))
		return c.NoContent(http.StatusInternalServerError)
	}

	ctx := context.Background()
	user, err := h.Repository.GetUserByName(ctx, traqUser.Name)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	res := &UserOverview{
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Admin:       user.Admin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
	return c.JSON(http.StatusOK, res)
}
