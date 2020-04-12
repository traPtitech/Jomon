// +build debug

package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/Jomon/model"
	storagePkg "github.com/traPtitech/Jomon/storage"
	"os"
)

//noinspection GoDuplicate
func NewService() Service {
	fmt.Printf("\n !!! THIS IS DEBUG BUILD !!!\n\n")

	dir := os.Getenv("UPLOAD_DIR")
	if dir == "" {
		dir = "./uploads"
	}
	local, err := storagePkg.NewLocalStorage(dir)
	if err != nil {
		panic(err)
	}

	s := Service{
		Administrators: model.NewAdministratorRepository(),
		Applications:   model.NewApplicationRepository(),
		Comments:       model.NewCommentRepository(),
		Images:         model.NewApplicationsImageRepository(&local),
		Users: &debugUserRepository{
			users: []model.User{
				{TrapId: "MyUser"},
				{TrapId: "AdminUser"},
				{TrapId: "NormalUser1"},
				{TrapId: "NormalUser2"},
			},
		},
		TraQAuth: model.NewTraQAuthRepository(""),
	}

	_ = s.Administrators.AddAdministrator("MyUser")
	_ = s.Administrators.AddAdministrator("AdminUser")

	return s
}

func EchoConfig(e *echo.Echo) {
	e.Use(middleware.Logger())
}

type debugUserRepository struct {
	users []model.User
}

func (d *debugUserRepository) GetUsers(token string) ([]model.User, error) {
	return d.users, nil
}

func (d *debugUserRepository) GetMyUser(token string) (model.User, error) {
	return d.users[0], nil
}

func (d *debugUserRepository) ExistsUser(token string, trapId string) (bool, error) {
	for _, user := range d.users {
		if trapId == user.TrapId {
			return true, nil
		}
	}
	return false, nil
}

func (s Service) AuthUser(c echo.Context) (echo.Context, error) {
	user, _ := s.Users.GetMyUser("")
	admins, _ := s.Administrators.GetAdministratorList()
	user.GiveIsUserAdmin(admins)

	c.Set(contextUserKey, user)
	c.Set(contextAccessTokenKey, "")

	return c, nil
}
