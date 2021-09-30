//go:generate /home/shota/go/bin/mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	storagePkg "github.com/traPtitech/Jomon/storage"
)

type Service interface {
	CreateFile(src io.Reader, name string, mimetype string, requestID uuid.UUID) (*File, error)
	GetAccessToken(code string, codeVerifier string) (AuthResponse, error)
	GetClientId() string
	GetMe(token string) (*User, error)
	StrToDate(str string) (time.Time, error)
	StrToTime(str string) (time.Time, error)
	WebhookEventHandler(c echo.Context, reqBody, resBody []byte)
}
type Services struct {
	Repository model.Repository
	Auth       Auth
	Storage    storagePkg.Storage
	Webhook    Webhook
}

func NewServices(repo model.Repository, storage storagePkg.Storage) (Service, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	traQClientID := os.Getenv("TRAQ_CLIENT_ID")
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	webhookChannelId := os.Getenv("WEBHOOK_CHANNEL_ID")
	id := os.Getenv("WEBHOOK_ID")
	return &Services{
		Repository: repo,
		Auth: Auth{
			ClientID: traQClientID,
		},
		Storage: storage,
		Webhook: Webhook{
			Secret:    webhookSecret,
			ChannelId: webhookChannelId,
			ID:        id,
		},
	}, nil
}
