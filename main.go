package main

import (
	"context"
	"os"

	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/model"
	"github.com/traPtitech/Jomon/internal/router"
	"github.com/traPtitech/Jomon/internal/service"
	"github.com/traPtitech/Jomon/internal/storage"
	"github.com/traPtitech/Jomon/internal/traq"
	"github.com/traPtitech/Jomon/internal/webhook"
	"go.uber.org/zap"
)

func main() {
	// Setup ent client
	client, err := model.Connect()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Setup storage
	var strg service.Storage
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
	repo := model.NewEntRepository(client)
	migrateOptions := []model.MigrateOption{}
	if migrationsDir := os.Getenv("MIGRATIONS_DIR"); migrationsDir != "" {
		migrateOptions = append(migrateOptions, model.MigrationsDir(migrationsDir))
	}
	if err := repo.MigrateApply(context.Background(), migrateOptions...); err != nil {
		panic(err)
	}
	// Setup webhook service
	ws, err := webhook.Load()
	if err != nil {
		panic(err)
	}
	// Setup OIDC client
	oidcClientConfig := traq.ClientConfig{
		ClientID:     os.Getenv("TRAQ_CLIENT_ID"),
		ClientSecret: os.Getenv("TRAQ_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("TRAQ_OIDC_CALLBACK_URL"),
		ProviderURL:  "https://q.trap.jp",
	}
	oidcClient, err := traq.LoadClient(context.Background(), oidcClientConfig)
	if err != nil {
		panic(err)
	}
	service := service.New(repo, strg, oidcClient)

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

	sessionName := os.Getenv("SESSION_NAME")
	if sessionName == "" {
		sessionName = "session"
	}
	handlers := router.Handlers{
		WebhookService: ws,
		Service:        service,
		SessionName:    sessionName,
	}

	server := handlers.NewServer(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	// Start server
	logger.Fatal("failed to start server", zap.Error(server.Start(":"+port)))
}
