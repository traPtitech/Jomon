package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/router"
	storagePkg "github.com/traPtitech/Jomon/storage"
	"net/http"
	"os"
)

func main() {
	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = model.Migrate()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/", genRootHandler(err == nil))

	var storage storagePkg.Storage
	if os.Getenv("OS_CONTAINER") != "" {
		// Swiftオブジェクトストレージ
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
		storage = &swift
	} else {
		// ローカルストレージ
		dir := os.Getenv("UPLOAD_DIR")
		if dir == "" {
			dir = "./uploads"
		}
		local, err := storagePkg.NewLocalStorage(dir)
		if err != nil {
			panic(err)
		}
		storage = &local
	}

	service := router.Service{
		Administrators: model.NewAdministratorRepository(),
		Applications:   model.NewApplicationRepository(),
		Comments:       model.NewCommentRepository(),
		Images:         model.NewApplicationsImageRepository(storage),
		Users:          model.NewUserRepository(),
	}

	router.SetRouting(e, service)
	e.Start(":1323")
}

func genRootHandler(b bool) func(echo.Context) error {
	return func(c echo.Context) error {
		if b {
			return c.String(http.StatusOK, "Succeeded in access db.")
		} else {
			return c.String(http.StatusOK, "Failed to access db.")
		}
	}
}
