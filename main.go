package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/router"
	"net/http"
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
	router.SetRouting(e)
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
