package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	DisplayName string    `json:"display_name"`
	Name        string    `json:"name"`
	Admin       bool      `json:"admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TraQUser struct {
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

func (s *Services) GetMe(token string) (*User, error) {
	req, err := http.NewRequest("GET", TraQBaseURL+"/users/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	body, err := sendReq(req)
	if err != nil {
		return nil, err
	}

	user := User{}
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func sendReq(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("StatusCode is not 200")
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
