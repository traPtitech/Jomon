package model

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type traqRepositoryMock struct {
	mock.Mock
	token string
}

var myUserId = "MyUserId"

func NewTraqRepositoryMock(token string) *traqRepositoryMock {
	m := new(traqRepositoryMock)
	m.token = token

	getUsersReq, err := http.NewRequest("GET", TraQBaseURL+"/users", nil)
	if err != nil {
		panic(err)
	}
	getUsersReq.Header.Set("Authorization", token)

	m.On("sendReq", getUsersReq).Return([]byte(fmt.Sprintf(`
	[
		{
			"name": "UserId"
		},
		{
			"name": "%s"
		}
	]`, myUserId)), nil)

	getMyUserReq, err := http.NewRequest("GET", TraQBaseURL+"/users/me", nil)
	if err != nil {
		panic(err)
	}
	getMyUserReq.Header.Set("Authorization", token)

	m.On("sendReq", getMyUserReq).Return([]byte(fmt.Sprintf(`
	{
		"name": "%s"
	}`, myUserId)), nil)

	return m
}

func (m *traqRepositoryMock) sendReq(req *http.Request) ([]byte, error) {
	ret := m.Called(req)
	return ret.Get(0).([]byte), ret.Error(1)
}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	token := "Token"

	userRepo := &userRepository{
		traqRepository: NewTraqRepositoryMock(token),
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		_, err := userRepo.GetUsers(token)
		asr.NoError(err)
	})
}

func TestGetMyUser(t *testing.T) {
	t.Parallel()

	token := "Token"

	userRepo := &userRepository{
		traqRepository: NewTraqRepositoryMock(token),
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		user, err := userRepo.GetMyUser(token)
		asr.NoError(err)
		asr.Equal(user.TrapId, myUserId)
	})
}

func TestExistsUser(t *testing.T) {
	t.Parallel()

	token := "Token"

	userRepo := &userRepository{
		traqRepository: NewTraqRepositoryMock(token),
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		exist, err := userRepo.ExistsUser(token, myUserId)
		asr.NoError(err)
		asr.True(exist)
	})

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		exist, err := userRepo.ExistsUser(token, "notExistId")
		asr.NoError(err)
		asr.False(exist)
	})
}
