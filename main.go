package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/Jomon/router"
	"go.uber.org/zap"
)

func main() {
	// setup server
	logger, _ := zap.NewDevelopment()
	// start server
	e := echo.New()
	e.Debug = (os.Getenv("IS_DEBUG_MODE") != "")
	e.Use(router.AccessLoggingMiddleware(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	router.SetRouting(e, service.NewService())
	e.Start(":1323")
}
