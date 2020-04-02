package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/router"
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

	router.SetRouting(e, router.NewService())
	e.Start(":1323")
}
