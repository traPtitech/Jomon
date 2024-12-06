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
	Group     *Group    `json:"group"`
}

type CommentApplication struct {
	ID      uuid.UUID `json:"id"`
	User    uuid.UUID `json:"user"`
	Comment string    `json:"comment"`
}

type TransactionPostRequestApplication struct {
	ID     uuid.UUID `json:"id"`
	Amount int       `json:"amount"`
	Target string    `json:"target"`
	Tags   []*Tag    `json:"tags"`
	Group  *Group    `json:"group"`
}

type TransactionPutRequestApplication struct {
	ID     uuid.UUID `json:"id"`
	Amount int       `json:"amount"`
	Target string    `json:"target"`
	Tags   []*Tag    `json:"tags"`
	Group  *Group    `json:"group"`
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

type Group struct {
	Name string `json:"name"`
}

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Admin       bool      `json:"admin"`
}

type Webhook struct {
	Secret    string
	ChannelId string
	ID        string
}

func WebhookEventHandler(c echo.Context, reqBody, resBody []byte) {
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	webhookChannelId := os.Getenv("WEBHOOK_CHANNEL_ID")
	webhookId := os.Getenv("WEBHOOK_ID")
	var message string

	if strings.Contains(c.Request().URL.Path, "/api/requests") {
		if strings.Contains(c.Request().URL.Path, "/comments") {
			resApp := new(CommentApplication)
			err := json.Unmarshal(resBody, resApp)
			if err != nil {
				return
			}
			splitedPath := strings.Split(c.Request().URL.Path, "/")

			message += fmt.Sprintf(
				"## :comment:[申請](%s/requests/%s)",
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

			if resApp.Group != nil {
				message += fmt.Sprintf("- 請求先グループ: %s\n", resApp.Group.Name)
			}

			if len(resApp.Tags) != 0 {
				tags := lo.Map(resApp.Tags, func(tag *Tag, _ int) string {
					return tag.Name
				})
				message += fmt.Sprintf("- タグ: %s", strings.Join(tags, ", "))
			}
			message += "\n\n"
			message += resApp.Content + "\n"
		}
	}
	_ = RequestWebhook(message, webhookSecret, webhookChannelId, webhookId, 1)
}

func WebhookPostEventHandler(c echo.Context, reqBody, resBody []byte) {
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	webhookChannelId := os.Getenv("WEBHOOK_CHANNEL_ID")
	webhookId := os.Getenv("WEBHOOK_ID")
	var message string
	var resApps []TransactionPostRequestApplication
	err := json.Unmarshal(resBody, &resApps)
	if err != nil || len(resApps) < 1 {
		return
	}
	message += fmt.Sprintf(
		"## :scroll:[入出金記録](%s/transactions/%s)が新規作成されました\n",
		"https://jomon.trap.jp",
		resApps[0].ID)
	targets := lo.Map(
		resApps, func(resApp TransactionPostRequestApplication, _ int) string {
			return resApp.Target
		})
	if resApps[0].Amount < 0 {
		message += fmt.Sprintf(
			"- %sへの支払い\n    - 支払い金額: 計%d円(一人当たりへの支払い金額: %d円)\n",
			strings.Join(targets, " "),
			-len(resApps)*resApps[0].Amount,
			-resApps[0].Amount)
	} else {
		message += fmt.Sprintf(
			"- %sからの振込\n    - 受け取り金額: 計%d円(一人当たりからの受け取り金額: %d円)\n",
			strings.Join(targets, " "),
			len(resApps)*resApps[0].Amount,
			resApps[0].Amount)
	}
	if resApps[0].Group != nil {
		message += fmt.Sprintf("- 関連するグループ: %s\n", resApps[0].Group.Name)
	}
	if len(resApps[0].Tags) != 0 {
		tags := lo.Map(resApps[0].Tags, func(tag *Tag, _ int) string {
			return tag.Name
		})
		message += fmt.Sprintf("- タグ: %s", strings.Join(tags, ", "))
	}

	_ = RequestWebhook(message, webhookSecret, webhookChannelId, webhookId, 1)
}

func WebhookPutEventHandler(c echo.Context, reqBody, resBody []byte) {
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	webhookChannelId := os.Getenv("WEBHOOK_CHANNEL_ID")
	webhookId := os.Getenv("WEBHOOK_ID")
	var message string
	var resApp TransactionPutRequestApplication
	err := json.Unmarshal(resBody, &resApp)
	if err != nil {
		return
	}
	message += fmt.Sprintf(
		"## :scroll:[入出金記録](%s/transactions/%s)が修正されました\n",
		"https://jomon.trap.jp",
		resApp.ID)
	if resApp.Amount < 0 {
		message += fmt.Sprintf(
			"- `%s`への支払い\n    - 支払い金額: %d円\n",
			resApp.Target,
			-resApp.Amount)
	} else {
		message += fmt.Sprintf(
			"- `%s`からの振込\n    - 受け取り金額: %d円\n",
			resApp.Target,
			resApp.Amount)
	}
	if resApp.Group != nil {
		message += fmt.Sprintf("- 関連するグループ: %s\n", resApp.Group.Name)
	}
	if len(resApp.Tags) != 0 {
		tags := lo.Map(resApp.Tags, func(tag *Tag, _ int) string {
			return tag.Name
		})
		message += fmt.Sprintf("- タグ: %s", strings.Join(tags, ", "))
	}

	_ = RequestWebhook(message, webhookSecret, webhookChannelId, webhookId, 1)
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
