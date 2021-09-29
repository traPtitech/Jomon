//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"os"
)

type Service interface {
	GetAccessToken(code string, codeVerifier string) (AuthResponse, error)
	GetClientId() string
	GetMe(token string) (*User, error)
}
type Services struct {
	Auth Auth
}

func NewServices() (*Services, error) {
	traQClientID := os.Getenv("TRAQ_CLIENT_ID")
	/*
		webhookSecret := os.Getenv("WEBHOOK_SECRET")
		webhookChannelId := os.Getenv("WEBHOOK_CHANNEL_ID")
		webhookId := os.Getenv("WEBHOOK_ID")
	*/
	return &Services{
		Auth: Auth{
			ClientID: traQClientID,
		},
	}, nil
}
