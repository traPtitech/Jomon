package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Admin       bool   `json:"admin"`
}

type UserOverview struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Admin       bool      `json:"admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserResponse struct {
	Users []*UserOverview `json:"users"`
}

func (h *Handlers) GetUsers(c echo.Context) error {
	ctx := context.Background()
	users, err := h.Repository.GetUsers(ctx)
	if err != nil {
		internalServerError(err)
	}
	res := []*UserOverview{}
	for _, user := range users {
		res = append(res, &UserOverview{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Admin:       user.Admin,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}
	return c.JSON(http.StatusOK, &UserResponse{res})
}

func (h *Handlers) PutUsers(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) GetMe(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
