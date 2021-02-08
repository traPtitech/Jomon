package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// TrapUser traP User struct
type TrapUser struct {
	TrapID  string `gorm:"type:varchar(32);not null;" json:"trap_id"`
	IsAdmin bool   `gorm:"-" json:"is_admin"`
}

// GiveIsUserAdmin check whether trapuser is admin or not
func (user *TrapUser) GiveIsUserAdmin(admins []string) {
	if user == nil {
		return
	}

	user.IsAdmin = false

	for _, admin := range admins {
		if user.TrapID == admin {
			user.IsAdmin = true
			break
		}
	}
}

// UserRepository Repo of User
type UserRepository interface {
	GetUsers(token string) ([]TrapUser, error)
	GetMyUser(token string) (TrapUser, error)
	ExistsUser(token string, trapID string) (bool, error)
}

type userRepository struct {
	traqRepository TraqRepository
}

// NewUserRepository Make UserRepository
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

// TraQBaseURL traQURL
const TraQBaseURL = "https://q.trap.jp/api/v3"

// TraqRepository Repo of traQ
type TraqRepository interface {
	sendReq(req *http.Request) ([]byte, error)
}

type traqRepository struct{}

// NewTraqRepository Make TraqRepository
func NewTraqRepository() TraqRepository {
	return &traqRepository{}
}

func (*traqRepository) sendReq(req *http.Request) ([]byte, error) {
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

func (repo *userRepository) GetUsers(token string) ([]TrapUser, error) {
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

	users := []TrapUser{}
	for _, traqUser := range traqUsers {
		if traqUser.Bot {
			continue
		}

		user := TrapUser{
			TrapID: traqUser.Name,
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepository) GetMyUser(token string) (TrapUser, error) {
	req, err := http.NewRequest("GET", TraQBaseURL+"/users/me", nil)
	if err != nil {
		return TrapUser{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	body, err := repo.traqRepository.sendReq(req)
	if err != nil {
		return TrapUser{}, err
	}

	traqUser := traqUser{}
	if err = json.Unmarshal(body, &traqUser); err != nil {
		return TrapUser{}, err
	}

	return TrapUser{
		TrapID: traqUser.Name,
	}, nil
}

func (repo *userRepository) ExistsUser(token string, trapID string) (bool, error) {
	users, err := repo.GetUsers(token)
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if trapID == user.TrapID {
			return true, nil
		}
	}

	return false, nil
}
