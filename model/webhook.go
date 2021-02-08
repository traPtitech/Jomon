package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

// WebhookRequest Response of BodyDump
type WebhookRequest struct {
	RequestID     uuid.UUID            `json:"request_id"`
	CurrentDetail WebhookRequestDetail `json:"current_detail"`
}

// WebhookRequestDetail Detail of Response of BodyDump
type WebhookRequestDetail struct {
	UpdateUser TrapUser  `json:"update_user"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Remarks    string    `json:"remarks"`
	Amount     int       `json:"amount"`
	PaidAt     time.Time `json:"paid_at"`
}

// WebhookRepository Repo of Webhook
type WebhookRepository interface {
	WebhookEventHandler(c echo.Context, reqBody, resBody []byte)
}

type webhookRepository struct {
	secret    string
	channelID string
	id        string
}

// NewWebhookRepository Make WebhookRepository
func NewWebhookRepository(secret string, channelID string, id string) WebhookRepository {
	return &webhookRepository{
		secret:    secret,
		channelID: channelID,
		id:        id,
	}
}

func (repo *webhookRepository) WebhookEventHandler(c echo.Context, reqBody, resBody []byte) {
	resApp := new(WebhookRequest)
	err := json.Unmarshal(resBody, resApp)
	if err != nil {
		return
	}
	var content string

	content = "## 申請書が作成されました" + "\n"
	content += fmt.Sprintf("### [%s](%s/applications/%s)", resApp.CurrentDetail.Title, "https://jomon.trap.jp", resApp.RequestID) + "\n"
	content += fmt.Sprintf("- 支払日: %s", resApp.CurrentDetail.PaidAt.Format("2006-01-02")) + "\n"
	content += fmt.Sprintf("- 支払金額: %s", strconv.Itoa(resApp.CurrentDetail.Amount)) + "\n"
	content += "\n"
	content += resApp.CurrentDetail.Remarks

	_ = RequestWebhook(content, repo.secret, repo.channelID, repo.id, 1)
}

// RequestWebhook send request to traQ api
func RequestWebhook(message, secret, channelID, webhookID string, embed int) error {
	u, err := url.Parse("https://q.trap.jp/api/v3/webhooks")
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, webhookID)
	query := u.Query()
	query.Set("embed", strconv.Itoa(embed))
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(message))
	if err != nil {
		return err
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	req.Header.Set("X-TRAQ-Signature", calcSignature(message, secret))
	if channelID != "" {
		req.Header.Set("X-TRAQ-Channel-Id", channelID)
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

func calcSignature(message, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
