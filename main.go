package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/srinathgs/mysqlstore"
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

	sessionStore, err := mysqlstore.NewMySQLStoreFromConnection(db.DB(), "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(session.Middleware(sessionStore))
	router.SetRouting(e, router.NewService())
	router.EchoConfig(e)
	e.Start(":1323")
}
