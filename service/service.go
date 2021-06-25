package service

import (
	"io"
	"os"
	"time"

	"github.com/traPtitech/Jomon/model"
	storagePkg "github.com/traPtitech/Jomon/storage"
)

type Service interface {
	CreateFile(src io.Reader, mimetype string) (File, error)
	GetAccessToken(code string, codeVerifier string) (AuthResponse, error)
	GetClientId() string
	StrToDate(str string) (time.Time, error)
	StrToTime(str string) (time.Time, error)
}
type Services struct {
	Repository model.Repository
	Auth       Auth
	Storage    storagePkg.Storage
}

func NewServices(repo model.Repository, storage storagePkg.Storage) (*Services, error) {
	traQClientID := os.Getenv("TRAQ_CLIENT_ID")
	/*
		webhookSecret := os.Getenv("WEBHOOK_SECRET")
		webhookChannelId := os.Getenv("WEBHOOK_CHANNEL_ID")
		webhookId := os.Getenv("WEBHOOK_ID")
	*/
	return &Services{
		Repository: repo,
		Auth: Auth{
			ClientID: traQClientID,
		},
		Storage: storage,
	}, nil
}
