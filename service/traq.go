package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type Auth struct {
	ClientID string
}

type Authority struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

const TraQBaseURL = "https://q.trap.jp/api/v3"

var JomonClientID = os.Getenv("TRAQ_CLIENT_ID")

func RequestAccessToken(code, codeVerifier string) (Authority, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", JomonClientID)
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)

	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", TraQBaseURL+"/oauth2/token", reqBody)
	if err != nil {
		return Authority{}, err
	}
	req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return Authority{}, err
	} else if res.StatusCode != http.StatusOK {
		return Authority{}, fmt.Errorf("failed to acquire access token")
	}

	var authRes Authority
	err = json.NewDecoder(res.Body).Decode(&authRes)
	if err != nil {
		return Authority{}, err
	}

	return authRes, nil
}

type TraQUser struct {
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

func FetchTraQUserInfo(token string) (TraQUser, error) {
	req, err := http.NewRequest("GET", TraQBaseURL+"/users/me", nil)
	if err != nil {
		return TraQUser{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TraQUser{}, err
	} else if res.StatusCode != http.StatusOK {
		return TraQUser{}, fmt.Errorf("failed to fetch user info")
	}

	var user TraQUser
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return TraQUser{}, err
	}

	return user, nil
}
