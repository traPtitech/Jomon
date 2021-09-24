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
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WebhookApplication struct {
	ApplicationID uuid.UUID                `json:"application_id"`
	CurrentDetail WebhookApplicationDetail `json:"current_detail"`
}

type WebhookApplicationDetail struct {
	UpdateUser User   `json:"update_user"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Remarks    string `json:"remarks"`
	Amount     int    `json:"amount"`
	PaidAt     PaidAt `json:"paid_at"`
}

type PaidAt struct {
	PaidAt time.Time `gorm:"type:date;not null"`
}

type WebhookRepository interface {
	WebhookEventHandler(c echo.Context, reqBody, resBody []byte)
}

type webhookRepository struct {
	secret    string
	channelId string
	id        string
}

func (repo *webhookRepository) WebhookEventHandler(c echo.Context, reqBody, resBody []byte) {
	resApp := new(WebhookApplication)
	err := json.Unmarshal(resBody, resApp)
	if err != nil {
		return
	}
	var content string

	content += "## 申請書が作成されました" + "\n"
	content += fmt.Sprintf("### [%s](%s/applications/%s)", resApp.CurrentDetail.Title, "https://jomon.trap.jp", resApp.ApplicationID) + "\n"
	content += fmt.Sprintf("- 支払日: %s", resApp.CurrentDetail.PaidAt.PaidAt.Format("2006-01-02")) + "\n"
	content += fmt.Sprintf("- 支払金額: %s", strconv.Itoa(resApp.CurrentDetail.Amount)) + "\n"
	content += "\n"
	content += resApp.CurrentDetail.Remarks

	_ = RequestWebhook(content, repo.secret, repo.channelId, repo.id, 1)

}

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
