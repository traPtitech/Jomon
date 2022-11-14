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

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type RequestApplication struct {
	ID        uuid.UUID `json:"id"`
	CreatedBy uuid.UUID `json:"created_by"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Amount    int       `json:"amount"`
	Tags      []*Tags   `json:"tags"`
	Group     *Group    `json:"group"`
}

type CommentApplication struct {
	ID      uuid.UUID `json:"id"`
	User    uuid.UUID `json:"user"`
	Comment string    `json:"comment"`
}

type TransactionRequestApplication struct {
	ID     uuid.UUID `json:"id"`
	Amount int       `json:"amount"`
	Target string    `json:"target"`
	Tags   []*Tags   `json:"tags"`
	Group  *Group    `json:"group"`
}

type Tags struct {
	Name string `json:"name"`
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

			message += fmt.Sprintf("## :comment:[依頼](%s/requests/%s)", "https://jomon.trap.jp", splitedPath[3])
			message += "に対する"
			message += fmt.Sprintf("[コメント](%s/requests/%s/comments/%s)", "https://jomon.trap.jp", splitedPath[3], resApp.ID)
			message += "が作成されました" + "\n"

			sess, _ := session.Get("session", c)
			bodyUser, _ := sess.Values["user"].([]byte)
			user := new(User)
			_ = json.Unmarshal(bodyUser, user)

			message += fmt.Sprintf("- 作成者: @%s", user.Name)
			message += "\n" + "\n"
			message += resApp.Comment + "\n"
		} else {
			resApp := new(RequestApplication)
			err := json.Unmarshal(resBody, resApp)
			if err != nil {
				return
			}
			if c.Request().Method == http.MethodPost {
				message += "## :receipt:依頼が作成されました" + "\n"
			} else if c.Request().Method == http.MethodPut {
				message += "## :receipt:依頼が更新されました" + "\n"
		

			message += fmt.Sprintf("### [%s](%s/applications/%s)", resApp.Title, "https://jomon.trap.jp", resApp.ID) + "\n"
			message += fmt.Sprintf("- 支払金額: %s円", strconv.Itoa(resApp.Amount)) + "\n"

			if resApp.Group != nil {
				message += fmt.Sprintf("- 請求先グループ: %s", resApp.Group.Name) + "\n"
			}

			if resApp.Tags != nil {
				message += "- タグ: "
				for _, tag := range resApp.Tags {
					message += tag.Name + ", "
				}
				message = message[:len(message)-len(", ")]
			}
			message += "\n" + "\n"
			message += resApp.Content + "\n"
		}
		}
	} else if strings.Contains(c.Request().URL.Path, "/api/transactions") {
		resApp := new(TransactionRequestApplication)
		err := json.Unmarshal(resBody, resApp)
		if err != nil {
			return
		}
		if c.Request().Method == http.MethodPost {
			message += "## :scroll:"
			message += fmt.Sprintf("[入出金記録](%s/applications/%s)", "https://jomon.trap.jp", resApp.ID) + "が新規作成されました" + "\n"
		} else if c.Request().Method == http.MethodPut {
			message += "## :scroll:"
			message += fmt.Sprintf("[入出金記録](%s/applications/%s)", "https://jomon.trap.jp", resApp.ID) + "が修正されました" + "\n"
		}
		if resApp.Amount < 0 {
			if resApp.Group != nil {
				message += fmt.Sprintf("- %s", resApp.Target) + "への依頼" + "\n"
			}
			message += fmt.Sprintf("- 支払い金額: % -s円", strconv.Itoa(resApp.Amount)) + "\n"
		} else {
			if resApp.Group != nil {
				message += fmt.Sprintf("- %s", resApp.Target) + "からの振込" + "\n"
			}
			message += fmt.Sprintf("- 受け取り金額: %s円", strconv.Itoa(resApp.Amount)) + "\n"
		}
		message += fmt.Sprintf("- 関連するグループ: %s", resApp.Group.Name) + "\n"	
		if resApp.Tags != nil {
			message += "- タグ: "
			tags := make([]string, len(resApp.Tags))
			for i, tag := range resApp.Tags {
				tags[i] = tag.Name
			}

			message += strings.Join(tags, ",")
			message = message[:len(message)-len(", ")]
		}
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
