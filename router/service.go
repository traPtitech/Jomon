// +build !debug

package router

import (
	"github.com/traPtitech/Jomon/model"
	storagePkg "github.com/traPtitech/Jomon/storage"
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

	return Service{
		Administrators: model.NewAdministratorRepository(),
		Applications:   model.NewApplicationRepository(),
		Comments:       model.NewCommentRepository(),
		Images:         model.NewApplicationsImageRepository(&swift),
		Users:          model.NewUserRepository(),
	}
}
