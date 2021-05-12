package main

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/router"
	"github.com/traPtitech/Jomon/service"
	"go.uber.org/zap"
)

func main() {
	// setup database
	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}

	err = model.Migrate(db)
	if err != nil {
		panic(err)
	}
	repo := model.NewGormRepository(db)

	// setup service
	services, err := service.NewServices(repo)
	if err != nil {
		panic(err)
	}

	// setup server
	logger, _ := zap.NewDevelopment()
	handlers := router.Handlers{
		Repo:         repo,
		Service:      services,
		SessionName:  "session",
		SessionStore: sessions.NewCookieStore([]byte("session")),
	}
	e := echo.New()
	e.Debug = (os.Getenv("IS_DEBUG_MODE") != "")
	e.Use(handlers.AccessLoggingMiddleware(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(handlers.SessionStore))

	handlers.Setup(e)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	e.Start(":" + port)
}
