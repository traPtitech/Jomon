// +build !debug

package router

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	storagePkg "github.com/traPtitech/Jomon/storage"
	"net/http"
	"os"
)

func NewService() Service {
	swift, err := storagePkg.NewSwiftStorage(
		os.Getenv("OS_CONTAINER"),
		os.Getenv("OS_USERNAME"),
		os.Getenv("OS_PASSWORD"),
		os.Getenv("OS_TENANT_NAME"),
		os.Getenv("OS_TENANT_ID"),
		os.Getenv("OS_AUTH_URL"),
	)
	if err != nil {
		panic(err)
	}

	traQClientId := os.Getenv("TRAQ_CLIENT_ID")

	return Service{
		Administrators: model.NewAdministratorRepository(),
		Applications:   model.NewApplicationRepository(),
		Comments:       model.NewCommentRepository(),
		Images:         model.NewApplicationsImageRepository(&swift),
		Users:          model.NewUserRepository(),
		TraQAuth:       model.NewTraQAuthRepository(traQClientId),
	}
}

func (s Service) AuthUser(c echo.Context) (echo.Context, error) {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c, c.NoContent(http.StatusInternalServerError)
	}

	accTok := sess.Values["accessToken"]
	if accTok == nil {
		return c, c.NoContent(http.StatusUnauthorized)
	}

	user, ok := sess.Values["user_name"].(model.User)
	if !ok {
		user, err = s.Users.GetMyUser(accTok.(string))

		if err != nil {
			return c, c.NoContent(http.StatusInternalServerError)
		}
	}
	c.Set("user", user)

	return c, nil
}
