package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type Auth struct {
	ClientID string
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

const TraQBaseURL = "https://q.trap.jp/api/v3"

func (s *Services) GetAccessToken(code string, codeVerifier string) (AuthResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", s.Auth.ClientID)
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", TraQBaseURL+"/oauth2/token", reqBody)
	if err != nil {
		return AuthResponse{}, err
	}
	req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return AuthResponse{}, err
	} else if res.StatusCode != 200 {
		return AuthResponse{}, fmt.Errorf("failed to acquire access token")
	}

	var authRes AuthResponse
	err = json.NewDecoder(res.Body).Decode(&authRes)
	if err != nil {
		return AuthResponse{}, err
	}

	return authRes, nil
}

func (s *Services) GetClientId() string {
	return s.Auth.ClientID
}
