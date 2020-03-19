package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type User struct {
	TrapId  string `gorm:"type:varchar(32);not null;" json:"trap_id"`
	IsAdmin bool   `gorm:"-" json:"is_admin"`
}

func (user *User) GiveIsUserAdmin(admins []string) {
	if user == nil {
		return
	}

	user.IsAdmin = false

	for _, admin := range admins {
		if user.TrapId == admin {
			user.IsAdmin = true
			break
		}
	}
}

type UserRepository interface {
	GetUsers(token string) ([]User, error)
	GetMyUser(token string) (User, error)
	IsUserFound(token string, trapId string) (bool, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

type traqUser struct {
	Name string `json:"name"`
}

const baseURL = "https://q.trap.jp/api/1.0"

func sendReqTraq(req *http.Request) ([]byte, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("StatusCode is not 200")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (_ *userRepository) GetUsers(token string) ([]User, error) {
	req, err := http.NewRequest("GET", baseURL+"/users", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	body, err := sendReqTraq(req)
	if err != nil {
		return nil, err
	}

	traqUsers := []traqUser{}
	if err = json.Unmarshal(body, &traqUsers); err != nil {
		return nil, err
	}

	users := []User{}
	for _, traqUser := range traqUsers {
		user := User{
			TrapId: traqUser.Name,
		}
		users = append(users, user)
	}

	return users, nil
}

func (_ *userRepository) GetMyUser(token string) (User, error) {
	req, err := http.NewRequest("GET", baseURL+"/users/me", nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Set("Authorization", token)

	body, err := sendReqTraq(req)
	if err != nil {
		return User{}, err
	}

	traqUser := traqUser{}
	if err = json.Unmarshal(body, &traqUser); err != nil {
		return User{}, err
	}

	return User{
		TrapId: traqUser.Name,
	}, nil
}

func (repo *userRepository) IsUserFound(token string, trapId string) (bool, error) {
	users, err := repo.GetUsers(token)
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if trapId == user.TrapId {
			return true, nil
		}
	}

	return false, nil
}
