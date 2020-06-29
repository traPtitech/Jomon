// +build !debug

package router

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	storagePkg "github.com/traPtitech/Jomon/storage"
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

	gob.Register(model.User{})

	return Service{
		Administrators: model.NewAdministratorRepository(),
		Applications:   model.NewApplicationRepository(),
		Comments:       model.NewCommentRepository(),
		Images:         model.NewApplicationsImageRepository(&swift),
		Users:          model.NewUserRepository(),
		TraQAuth:       model.NewTraQAuthRepository(traQClientId),
	}
}

func EchoConfig(_ *echo.Echo) {}

func (s Service) AuthUser(c echo.Context) (echo.Context, error) {
	sess, err := session.Get(sessionKey, c)
	if err != nil {
		return nil, c.NoContent(http.StatusInternalServerError)
	}

	accTok, ok := sess.Values[sessionAccessTokenKey].(string)
	if !ok || accTok == "" {
		return nil, c.NoContent(http.StatusUnauthorized)
	}
	c.Set(contextAccessTokenKey, accTok)

	user, ok := sess.Values[sessionUserKey].(model.User)
	if !ok {
		user, err = s.Users.GetMyUser(accTok)
		sess.Values[sessionUserKey] = user
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return nil, c.NoContent(http.StatusInternalServerError)
		}

		if err != nil {
			return nil, c.NoContent(http.StatusInternalServerError)
		}
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return nil, c.NoContent(http.StatusInternalServerError)
	}
	user.GiveIsUserAdmin(admins)

	c.Set(contextUserKey, user)

	return c, nil
}
