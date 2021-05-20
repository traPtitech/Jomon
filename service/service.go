package service

import (
	"os"

	"github.com/traPtitech/Jomon/model"
)

type Service interface {
	GetAccessToken(code string, codeVerifier string) (AuthResponse, error)
	GetClientId() string
}
type Services struct {
	Repository model.Repository
	Auth       Auth
}

func NewServices(repo model.Repository) (Services, error) {
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
	}, nil
}
