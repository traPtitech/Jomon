package router

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/model"
	"go.uber.org/zap"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"display_name"`
	Admin       bool       `json:"admin"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (h Handlers) GetUsers(c echo.Context) error {
	users, err := h.Repository.GetUsers(c.Request().Context())
	if err != nil {
		h.Logger.Error("failed to get users from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(users, func(user *model.User, _ int) User {
		return User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Admin:       user.Admin,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			DeletedAt:   user.DeletedAt,
		}
	})

	return c.JSON(http.StatusOK, res)
}

type PutUserRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Admin       bool   `json:"admin"`
}

func (h Handlers) UpdateUserInfo(c echo.Context) error {
	var newUser PutUserRequest
	if err := c.Bind(&newUser); err != nil {
		h.Logger.Info("could not get user info from request", zap.Error(err))
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := h.Repository.GetUserByName(c.Request().Context(), newUser.Name)
	if err != nil {
		h.Logger.Error("failed to get user from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updated, err := h.Repository.UpdateUser(
		c.Request().Context(),
		user.ID, newUser.Name, newUser.DisplayName, newUser.Admin)
	if err != nil {
		h.Logger.Error("failed to update user in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, User{
		ID:          user.ID,
		Name:        updated.Name,
		DisplayName: updated.DisplayName,
		Admin:       updated.Admin,
	})
}

func (h Handlers) GetMe(c echo.Context) error {
	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		h.Logger.Error("failed to get session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	userInSession, ok := sess.Values[sessionUserKey].(User)
	if !ok {
		h.Logger.Error("failed to parse stored session as user info")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user info")
	}

	user, err := h.Repository.GetUserByID(c.Request().Context(), userInSession.ID)
	if err != nil {
		users, err := h.Repository.GetUsers(c.Request().Context())
		for _, user_i := range users {
			if user_i.ID == userInSession.ID {
				return c.JSON(http.StatusOK, user_i)
			}
		}
		if err != nil {
			h.Logger.Error("failed to get users from repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		// if user not found, create new user
		new_user, err := h.Repository.CreateUser(c.Request().Context(), userInSession.Name,
			userInSession.DisplayName, userInSession.Admin)
		if err != nil {
			h.Logger.Error("failed to create user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, new_user)
	}
	return c.JSON(http.StatusOK, user)
}
