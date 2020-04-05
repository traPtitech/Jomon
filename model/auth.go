package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type TraQAuthRepository interface {
	GetAccessToken(code string, codeVerifier string) (AuthResponse, error)
	GetClientId() string
}

type traQAuthRepository struct {
	clientId string
}

func NewTraQAuthRepository(clientId string) TraQAuthRepository {
	return &traQAuthRepository{clientId: clientId}
}

func (repo *traQAuthRepository) GetAccessToken(code string, codeVerifier string) (AuthResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", repo.clientId)
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", TraQBaseURL+"/oauth2/token", reqBody)
	if err != nil {
		return AuthResponse{}, err
	}
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

func (repo *traQAuthRepository) GetClientId() string {
	return repo.clientId
}
