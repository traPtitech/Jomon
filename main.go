package main

import (
	"context"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
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

	// Setup model repository
	repo := model.NewEntRepository(client)
	// Setup service
	services, err := service.NewServices(repo)
	if err != nil {
		panic(err)
	}

	// Setup server
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
		Logger:       logger,
		Service:      services,
		SessionName:  "session",
		SessionStore: sessions.NewCookieStore([]byte("session")),
	}

	e := echo.New()

	router.SetRouting(e, handlers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	// Start server
	e.Logger.Fatal("failed to start server", e.Start(":"+port))
}
