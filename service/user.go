package service

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	TrapID    string    `json:"trap_id"`
	Name      string    `json:"string"`
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TraQUser struct {
	Name string `json:"name"`
	Bot  bool   `json:"bot"`
}

func (s *Services) GetMe(token string) (*User, error) {
	req, err := http.NewRequest("GET", TraQBaseURL+"/users/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	body, err := s.sendReq(req)
	if err != nil {
		return nil, err
	}

	traqUser := TraQUser{}
	if err = json.Unmarshal(body, &traqUser); err != nil {
		return nil, err
	}

	ctx := context.Background()
	user, err := s.Repository.GetMe(ctx, traqUser.Name)

	return ConvertModelUserToServiceUser(user), nil
}

func (_ *Services) sendReq(req *http.Request) ([]byte, error) {
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
