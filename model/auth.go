package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthResponse Response of Auth
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// TraQAuthRepository Repo to traQAuth
type TraQAuthRepository interface {
	GetAccessToken(code string, codeVerifier string) (AuthResponse, error)
	GetClientID() string
}

type traQAuthRepository struct {
	clientID string
}

// NewTraQAuthRepository Make TraQAuthRepository
func NewTraQAuthRepository(clientID string) TraQAuthRepository {
	return &traQAuthRepository{clientID: clientID}
}

func (repo *traQAuthRepository) GetAccessToken(code string, codeVerifier string) (AuthResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", repo.clientID)
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

func (repo *traQAuthRepository) GetClientID() string {
	return repo.clientID
}
