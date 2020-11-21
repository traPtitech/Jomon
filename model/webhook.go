package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type WebhookRepository interface {
	WebhookEventHandler(c echo.Context, reqBody, resBody []byte)
}

type webhookRepository struct {
	secret    string
	channelId string
	id        string
}

func NewWebhookRepository(secret string, channelId string, id string) WebhookRepository {
	return &webhookRepository{
		secret:    secret,
		channelId: channelId,
		id:        id,
	}
}

func (repo *webhookRepository) WebhookEventHandler(c echo.Context, reqBody, resBody []byte) {
	resApp := new(Application)
	err := json.Unmarshal(resBody, resApp)
	if err != nil {
		return
	}
	var content string

	content = "## 申請書が作成されました" + "\n"
	content += fmt.Sprintf("### [%s](%s/applications/%s)", resApp.LatestApplicationsDetail.Title, "https://jomon.trap.jp", resApp.ID) + "\n"
	content += fmt.Sprintf("- 支払日: %s", resApp.LatestApplicationsDetail.PaidAt.PaidAt.Format("2006-01-02")) + "\n"
	content += fmt.Sprintf("- 支払金額: %s", strconv.Itoa(resApp.LatestApplicationsDetail.Amount)) + "\n"
	content += "\n"
	content += resApp.LatestApplicationsDetail.Remarks

	u, err := url.Parse("https://q.trap.jp/api/v3/webhooks")
	if err != nil {
		return
	}

	u.Path = path.Join(u.Path, repo.id)
	query := u.Query()
	query.Set("embed", strconv.Itoa(1))
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(content))
	if err != nil {
		return
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)
	req.Header.Set("X-TRAQ-Signature", calcSignature(content, repo.secret))
	if repo.channelId != "" {
		req.Header.Set("X-TRAQ-Channel-Id", repo.channelId)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode >= 400 {
		return
	}
}

func calcSignature(message, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
