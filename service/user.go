package service

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/model"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
	Name        string    `json:"name"`
	Admin       bool      `json:"admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	user, err := s.Repository.GetUserByName(ctx, traqUser.Name)
	if err != nil {
		return nil, err
	}

	return ConvertModelUserToServiceUser(*user), nil
}

func (*Services) sendReq(req *http.Request) ([]byte, error) {
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

func ConvertModelUserToServiceUser(user model.User) *User {
	return &User{
		ID:          user.ID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Admin:       user.Admin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
