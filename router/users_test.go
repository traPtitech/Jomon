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
	userId      string
	adminUserId string
}

func NewUserRepositoryMock(userId string, adminUserId string) *userRepositoryMock {
	m := new(userRepositoryMock)
	m.userId = userId
	m.adminUserId = adminUserId

	m.On("GetUsers", mock.Anything).Return([]model.User{
		{TrapId: "User1"},
		{TrapId: "User2"},
		{TrapId: "User3"},
		{TrapId: userId},
		{TrapId: adminUserId},
	}, nil)

	m.On("ExistsUser", mock.Anything, userId).Return(true, nil)
	m.On("ExistsUser", mock.Anything, userId).Return(true, nil)
	m.On("ExistsUser", mock.Anything, adminUserId).Return(true, nil)
	m.On("ExistsUser", mock.Anything, adminUserId).Return(true, nil)
	m.On("ExistsUser", mock.Anything, mock.Anything).Return(false, nil)

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

func (m *userRepositoryMock) ExistsUser(token string, trapId string) (bool, error) {
	ret := m.Called(token, trapId)
	return ret.Bool(0), ret.Error(1)
}

func (m *userRepositoryMock) SetNormalUser(c echo.Context) {
	c.Set(contextUserKey, model.User{TrapId: m.userId})
}

func (m *userRepositoryMock) SetAnotherNormalUser(c echo.Context) {
	c.Set(contextUserKey, model.User{
		TrapId:  "User1",
		IsAdmin: false,
	})
}

func (m *userRepositoryMock) SetAdminUser(c echo.Context) {
	c.Set(contextUserKey, model.User{
		TrapId:  m.adminUserId,
		IsAdmin: true,
	})
}

func TestGetUsers(t *testing.T) {
	userId := "UserId"
	adminUserId := "AdminUserId"

	adminRepMock := NewAdministratorRepositoryMock(adminUserId)

	userRepMock := NewUserRepositoryMock(userId, adminUserId)

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users")
		userRepMock.SetNormalUser(c)
		c.Set(contextAccessTokenKey, "")

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

	userRepMock := NewUserRepositoryMock(userId, adminUserId)

	service := Service{
		Administrators: adminRepMock,
		Users:          userRepMock,
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/me")
		userRepMock.SetNormalUser(c)

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

	userRepMock := NewUserRepositoryMock(userId, adminUserId)

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")
		userRepMock.SetNormalUser(c)
		c.Set(contextAccessTokenKey, "")

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

	userRepMock := NewUserRepositoryMock(userId, adminUserId)

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")
		userRepMock.SetAdminUser(c)
		c.Set(contextAccessTokenKey, "")

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")
		userRepMock.SetAdminUser(c)
		c.Set(contextAccessTokenKey, "")

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")
		userRepMock.SetNormalUser(c)
		c.Set(contextAccessTokenKey, "")

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")
		userRepMock.SetAdminUser(c)
		c.Set(contextAccessTokenKey, "")

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/admins")
		userRepMock.SetAdminUser(c)
		c.Set(contextAccessTokenKey, "")

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

		err = service.PutAdminUsers(c)
		asr.NoError(err)
		asr.Equal(http.StatusNotFound, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
}
