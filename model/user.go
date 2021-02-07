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
type TrapUser struct {
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
	ExistsUser(token string, trapId string) (bool, error)
}

type userRepository struct {
	traqRepository TraqRepository
}

func NewUserRepository() UserRepository {
	return &userRepository{
		traqRepository: NewTraqRepository(),
	}
}

// v3ではdefaultでsuspendedは取得しない
type traqUser struct {
	Name string `json:"name"`
	Bot  bool   `json:"bot"`
}

const TraQBaseURL = "https://q.trap.jp/api/v3"

type TraqRepository interface {
	sendReq(req *http.Request) ([]byte, error)
}

type traqRepository struct{}

func NewTraqRepository() TraqRepository {
	return &traqRepository{}
}

func (_ *traqRepository) sendReq(req *http.Request) ([]byte, error) {
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

func (repo *userRepository) GetUsers(token string) ([]User, error) {
	req, err := http.NewRequest("GET", TraQBaseURL+"/users", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	body, err := repo.traqRepository.sendReq(req)
	if err != nil {
		return nil, err
	}

	traqUsers := []traqUser{}
	if err = json.Unmarshal(body, &traqUsers); err != nil {
		return nil, err
	}

	users := []User{}
	for _, traqUser := range traqUsers {
		if traqUser.Bot {
			continue
		}

		user := User{
			TrapId: traqUser.Name,
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepository) GetMyUser(token string) (User, error) {
	req, err := http.NewRequest("GET", TraQBaseURL+"/users/me", nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	body, err := repo.traqRepository.sendReq(req)
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

func (repo *userRepository) ExistsUser(token string, trapId string) (bool, error) {
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
