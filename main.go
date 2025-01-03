package main

import (
	"os"

	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/router"
	"github.com/traPtitech/Jomon/storage"
	"go.uber.org/zap"
)

func main() {
	// Setup ent client
	client, err := model.SetupEntClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Setup storage
	var strg storage.Storage
	if os.Getenv("IS_DEBUG_MODE") != "" {
		strg, err = storage.NewLocalStorage(os.Getenv("UPLOAD_DIR"))
		if err != nil {
			panic(err)
		}
	} else {
		strg, err = storage.NewSwiftStorage(
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
	}
	// Setup repository
	repo := model.NewEntRepository(client, strg)

	// Setup server
	logMode := logging.ModeFromEnv("IS_DEBUG_MODE")
	logger, err := logging.Load(logMode)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()
	handlers := router.Handlers{
		Repository:  repo,
		Storage:     strg,
		SessionName: "session",
	}

	server := handlers.NewServer(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	// Start server
	logger.Fatal("failed to start server", zap.Error(server.Start(":"+port)))
}
