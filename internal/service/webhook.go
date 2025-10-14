package service

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type RequestApplication struct {
	ID        uuid.UUID `json:"id"`
	CreatedBy uuid.UUID `json:"created_by"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []*Tag    `json:"tags"`
	Targets   []*Target `json:"targets"`
}

type CommentApplication struct {
	ID      uuid.UUID `json:"id"`
	User    uuid.UUID `json:"user"`
	Comment string    `json:"comment"`
}

type Tag struct {
	Name string `json:"name"`
}

type Target struct {
	ID        uuid.UUID `json:"id"`
	Target    uuid.UUID `json:"target"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	DisplayName    string    `json:"display_name"`
	AccountManager bool      `json:"account_manager"`
}

type Webhook struct {
	Secret    string
	ChannelId string
	ID        string
}

type WebhookService struct {
	// WEBHOOK_SECRET
	secret string
	// WEBHOOK_CHANNEL_ID
	channelID string
	// WEBHOOK_ID
	webhookID string
}

func LoadWebhookService() (*WebhookService, error) {
	loadEnv := func(key string) (string, error) {
		value := os.Getenv(key)
		if value == "" {
			return "", fmt.Errorf("%s is not set", key)
		}
		return value, nil
	}
	secret, err := loadEnv("WEBHOOK_SECRET")
	if err != nil {
		return nil, err
	}
	// WEBHOOK_CHANNEL_ID can be empty
	channelID := os.Getenv("WEBHOOK_CHANNEL_ID")
	webhookID, err := loadEnv("WEBHOOK_ID")
	if err != nil {
		return nil, err
	}

	return NewWebhookService(secret, channelID, webhookID), nil
}

func NewWebhookService(secret, channelID, webhookID string) *WebhookService {
	return &WebhookService{
		secret:    secret,
		channelID: channelID,
		webhookID: webhookID,
	}
}

func (ws *WebhookService) WebhookApplicationsEventHandler(c echo.Context, reqBody, resBody []byte) {
	var message string

	if strings.Contains(c.Request().URL.Path, "/comments") {
		resApp := new(CommentApplication)
		err := json.Unmarshal(resBody, resApp)
		if err != nil {
			return
		}
		splitedPath := strings.Split(c.Request().URL.Path, "/")

		message += fmt.Sprintf(
			"## :comment:[申請](%s/applications/%s)",
			"https://jomon.trap.jp",
			splitedPath[3])
		message += "に対する"
		message += fmt.Sprintf(
			"[コメント](%s/requests/%s/comments/%s)",
			"https://jomon.trap.jp",
			splitedPath[3],
			resApp.ID)
		message += "が作成されました\n\n"
		message += resApp.Comment + "\n"
	} else {
		resApp := new(RequestApplication)
		err := json.Unmarshal(resBody, resApp)
		if err != nil {
			return
		}
		if c.Request().Method == http.MethodPost {
			message += "## :receipt:申請が作成されました\n"
		} else if c.Request().Method == http.MethodPut {
			message += "## :receipt:申請が更新されました\n"
		}

		message += fmt.Sprintf(
			"### [%s](%s/applications/%s)\n",
			resApp.Title,
			"https://jomon.trap.jp",
			resApp.ID)

		amount := lo.Reduce(resApp.Targets, func(amo int, target *Target, _ int) int {
			return amo + target.Amount
		}, 0)
		message += fmt.Sprintf("- 支払金額: %d円\n", amount)

		if len(resApp.Tags) != 0 {
			tags := lo.Map(resApp.Tags, func(tag *Tag, _ int) string {
				return tag.Name
			})
			message += fmt.Sprintf("- タグ: %s", strings.Join(tags, ", "))
		}
		message += "\n\n"
		message += resApp.Content + "\n"
	}
	_ = ws.RequestWebhook(message, 1)
}

func (ws *WebhookService) RequestWebhook(message string, embed int) error {
	u, err := url.Parse("https://q.trap.jp/api/v3/webhooks")
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, ws.webhookID)
	query := u.Query()
	query.Set("embed", strconv.Itoa(embed))
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(message))
	if err != nil {
		return err
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	req.Header.Set("X-TRAQ-Signature", ws.calcSignature(message))
	if ws.channelID != "" {
		req.Header.Set("X-TRAQ-Channel-Id", ws.channelID)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return errors.New(http.StatusText(res.StatusCode))
	}
	return nil
}

func (ws *WebhookService) calcSignature(message string) string {
	mac := hmac.New(sha1.New, []byte(ws.secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
