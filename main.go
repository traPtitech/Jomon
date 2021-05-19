package main

import (
	"context"
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
	// Setup ent client
	client, err := model.SetupEntClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}

	// setup model repository
	repo := model.NewEntRepository(client)
	// setup service
	services, err := service.NewServices(repo)
	if err != nil {
		panic(err)
	}

	// setup server
	var logger *zap.Logger
	if os.Getenv("IS_DEBUG_MODE") != "" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	handlers := router.Handlers{
		Repository:   repo,
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

	router.SetRouting(e, handlers)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	e.Start(":" + port)

}
