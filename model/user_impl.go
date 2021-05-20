package model

import (
	"encoding/json"
	"net/http"
)

func (repo *EntRepository) GetMyUser(token string) (User, error) {
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
