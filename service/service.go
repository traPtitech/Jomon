package service

import (
	"io"
	"os"

	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/storage"
)

type Service interface {
	CreateFile(src io.Reader, mimetype string) (File, error)
	GetAccessToken(code string, codeVerifier string) (AuthResponse, error)
	GetClientId() string
}
type Services struct {
	Repository   model.Repository
	Auth         Auth
	SwiftStorage storage.Swift
}

func NewServices(repo model.Repository) (Services, error) {
	swift, err := storage.NewSwiftStorage(
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
	traQClientID := os.Getenv("TRAQ_CLIENT_ID")
	/*
		webhookSecret := os.Getenv("WEBHOOK_SECRET")
		webhookChannelId := os.Getenv("WEBHOOK_CHANNEL_ID")
		webhookId := os.Getenv("WEBHOOK_ID")
	*/
	return Services{
		Repository: repo,
		Auth: Auth{
			ClientID: traQClientID,
		},
		SwiftStorage: swift,
	}, nil
}
