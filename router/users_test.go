package router

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/traPtitech/Jomon/model"
)

type userRepositoryMock struct {
	mock.Mock
	asr        *assert.Assertions
	token      string
	adminToken string
}

func NewUserRepositoryMock(t *testing.T, userId string, adminUserId string) *userRepositoryMock {
	m := new(userRepositoryMock)
	m.asr = assert.New(t)
	m.token = "Token"
	m.adminToken = "AdminToken"

	m.On("GetUsers", m.token).Return([]model.User{
		{TrapId: "User1"},
		{TrapId: "User2"},
		{TrapId: "User3"},
		{TrapId: userId},
		{TrapId: adminUserId},
	}, nil)

	m.On("GetMyUser", m.token).Return(model.User{TrapId: userId}, nil)
	m.On("GetMyUser", m.adminToken).Return(model.User{TrapId: adminUserId}, nil)

	m.On("IsUserExist", m.token, userId).Return(true, nil)
	m.On("IsUserExist", m.adminToken, userId).Return(true, nil)
	m.On("IsUserExist", m.token, adminUserId).Return(true, nil)
	m.On("IsUserExist", m.adminToken, adminUserId).Return(true, nil)
	m.On("IsUserExist", mock.Anything, mock.Anything).Return(false, nil)

	return m
}

func (m *userRepositoryMock) GetUsers(token string) ([]model.User, error) {
	ret := m.Called(token)
	return ret.Get(0).([]model.User), ret.Error(1)
}

func (m *userRepositoryMock) GetMyUser(token string) (model.User, error) {
	ret := m.Called(token)
	return ret.Get(0).(model.User), ret.Error(1)
}

func (m *userRepositoryMock) IsUserExist(token string, trapId string) (bool, error) {
	ret := m.Called(token, trapId)
	return ret.Bool(0), ret.Error(1)
}

func TestGetUsers(t *testing.T) {
	userId := "UserId"
	adminUserId := "AdminUserId"

	adminRepMock := NewAdministratorRepositoryMock(adminUserId)

	userRepMock := NewUserRepositoryMock(t, userId, adminUserId)

	service := Service{
		Administrators: adminRepMock,
		Users:          userRepMock,
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.GetUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
}

func TestGetMyUser(t *testing.T) {
	userId := "UserId"
	adminUserId := "AdminUserId"

	adminRepMock := NewAdministratorRepositoryMock(adminUserId)

	userRepMock := NewUserRepositoryMock(t, userId, adminUserId)

	service := Service{
		Administrators: adminRepMock,
		Users:          userRepMock,
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/me")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.GetMyUser(c)
		asr.NoError(err)
		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
}

func TestGetAdminUsers(t *testing.T) {
	userId := "UserId"
	adminUserId := "AdminUserId"

	adminRepMock := NewAdministratorRepositoryMock(adminUserId)

	userRepMock := NewUserRepositoryMock(t, userId, adminUserId)

	service := Service{
		Administrators: adminRepMock,
		Users:          userRepMock,
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/api/users/admins", nil)
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.GetAdminUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
}

func TestPutAdminUsers(t *testing.T) {
	userId := "UserId"
	adminUserId := "AdminUserId"

	userRepMock := NewUserRepositoryMock(t, userId, adminUserId)

	adminRepMock := NewAdministratorRepositoryMock(adminUserId)

	adminRepMock.On("AddAdministrator", userId).Return(nil)
	adminRepMock.On("RemoveAdministrator", adminUserId).Return(nil)

	service := Service{
		Administrators: adminRepMock,
		Users:          userRepMock,
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"trap_id":"%s",
			"to_admin":true
		}
		`, userId)

		req := httptest.NewRequest(http.MethodPut, "/api/users/admins", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.adminToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutAdminUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"trap_id":"%s",
			"to_admin":false
		}
		`, adminUserId)

		req := httptest.NewRequest(http.MethodPut, "/api/users/admins", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.adminToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutAdminUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"trap_id":"%s",
			"to_admin":true
		}
		`, userId)

		req := httptest.NewRequest(http.MethodPut, "/api/users/admins", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutAdminUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusForbidden, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"trap_id":"%s",
			"to_admin":false
		}
		`, userId)

		req := httptest.NewRequest(http.MethodPut, "/api/users/admins", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.adminToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutAdminUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusConflict, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"trap_id":"%s",
			"to_admin":true
		}
		`, "notExistUser")

		req := httptest.NewRequest(http.MethodPut, "/api/users/admins", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.adminToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutAdminUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusNotFound, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
}
