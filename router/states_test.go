package router

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/traPtitech/Jomon/model"
)

func (m *commentRepositoryMock) GetStatesLog(applicationId uuid.UUID, statesLogId int) (model.StatesLog, error) {
	ret := m.Called(applicationId, statesLogId)
	return ret.Get(0).(model.StatesLog), ret.Error(1)
}

func (m *commentRepositoryMock) PutStatesLog(applicationId uuid.UUID, statesLogId int, toState model.StateType, reason string) (model.StatesLog, error) {
	ret := m.Called(applicationId, statesLogId, toState, reason)
	return ret.Get(0).(model.StatesLog), ret.Error(1)
}

func TestPutState(t *testing.T) {
	appRepMock := NewApplicationRepositoryMock(t)

	title := "夏コミの交通費をお願いします。"
	remarks := "〇〇駅から〇〇駅への移動"
	amount := 1000
	paidAt := time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	stateReason := "This is reason."

	appRepMock.On("GetApplication", id, mock.Anything).Return(GenerateApplication(id, "User2", model.ApplicationType{Type: model.Contest}, title, remarks, amount, paidAt), nil)
	appRepMock.On("GetApplication", mock.Anything, mock.Anything).Return(model.Application{}, gorm.ErrRecordNotFound)
	appRepMock.On("UpdateStatesLog", id, mock.Anything, mock.Anything, mock.Anything).Return(model.StatesLog{}, nil)

	adminRepMock := NewAdministratorRepositoryMock("AdminUserId")
	adminRepMock.On("IsAdmin", "UserId").Return(false, nil)
	adminRepMock.On("IsAdmin", "AnotherId").Return(false, nil)
	adminRepMock.On("IsAdmin", "AdminUserId").Return(true, nil)

	toStateAccepted, err := model.StateType{Type: model.Accepted}.MarshalJSON()
	if err != nil {
		panic(err)
	}
	toStateFixRequired, err := model.StateType{Type: model.FixRequired}.MarshalJSON()
	if err != nil {
		panic(err)
	}

	userRepMock := NewUserRepositoryMock("UserId", "AdminUserId")

	service := Service{
		Administrators: adminRepMock,
		Applications:   appRepMock,
		Users:          userRepMock,
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()
		body := fmt.Sprintf(`
		{
			"to_state": %s,
			"reason": "%s"
		}
		`, string(toStateAccepted), stateReason)
		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/states", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/states")
		c.SetParamNames("applicationId")
		c.SetParamValues(id.String())
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

		err = service.PutStates(c)
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
			"to_state": %s,
			"reason": ""
		}
		`, string(toStateFixRequired))

		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/states", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/states")
		c.SetParamNames("applicationId")
		c.SetParamValues(id.String())
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

		err = service.PutStates(c)
		asr.NoError(err)
		asr.Equal(http.StatusBadRequest, rec.Code)
	})
}

func TestPutRepaidStates(t *testing.T) {
	appRepMock := NewApplicationRepositoryMock(t)

	title := "夏コミの交通費をお願いします。"
	remarks := "〇〇駅から〇〇駅への移動"
	amount := 1000
	paidAt := time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	userId := "UserId"

	stateReason := "This is reason."

	dt := time.Now().Format("2006-01-02")

	appRepMock.On("GetApplication", id, mock.Anything).Return(GenerateApplicationStatesLogAccepted(id, "User2", model.ApplicationType{Type: model.Contest}, title, remarks, amount, paidAt), nil)
	appRepMock.On("GetApplication", mock.Anything, mock.Anything).Return(model.Application{}, gorm.ErrRecordNotFound)
	appRepMock.On("UpdateStatesLog", id, mock.Anything, mock.Anything, dt).Return(model.StatesLog{}, nil)
	appRepMock.On("UpdateRepayUser", id, mock.Anything, mock.Anything, mock.Anything).Return(model.RepayUser{}, true, nil)

	adminRepMock := NewAdministratorRepositoryMock("AdminUserId")
	adminRepMock.On("IsAdmin", userId).Return(true, nil)

	toStateAccepted, err := model.StateType{Type: model.Accepted}.MarshalJSON()
	if err != nil {
		panic(err)
	}

	userRepMock := NewUserRepositoryMock("UserId", "AdminUserId")

	service := Service{
		Administrators: adminRepMock,
		Applications:   appRepMock,
		Users:          userRepMock,
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()
		body := fmt.Sprintf(`
		{
			"to_state": %s,
			"reason": "%s"
		}
		`, string(toStateAccepted), stateReason)
		body2 := fmt.Sprintf(`
		{
			"repaid_at": "%s"
		}
		`, dt)
		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/states", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/states")
		c.SetParamNames("applicationId")
		c.SetParamValues(id.String())
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
		err = service.PutStates(c)
		if err != nil {
			panic(err)
		}

		req = httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/states/repaid/"+userId, strings.NewReader(body2))
		req.Header.Set(echo.HeaderContentType, "application/json")
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/states/repaid/:repaidToId")
		c.SetParamNames("applicationId", "repaidToId")
		c.SetParamValues(id.String(), userId)
		userRepMock.SetNormalUser(c)

		route, pathParam, err = router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput = &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		err = service.PutRepaidStates(c)
		asr.NoError(err)
		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body2 := fmt.Sprintf(`
		{
			"repaid_at": "%s"
		}
		`, dt)

		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/states/repaid/"+userId, strings.NewReader(body2))
		req.Header.Set(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/states/repaid/:repaidToId")
		c.SetParamNames("applicationId", "repaidToId")
		c.SetParamValues(id.String(), userId)
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

		err = service.PutRepaidStates(c)
		if err != nil {
			panic(err)
		}

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)

		asr = assert.New(t)
		e = echo.New()
		ctx = context.TODO()

		req = httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/states/repaid/"+userId, strings.NewReader(body2))
		req.Header.Set(echo.HeaderContentType, "application/json")
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/states/repaid/:repaidToId")
		c.SetParamNames("applicationId")
		c.SetParamValues(id.String())
		userRepMock.SetNormalUser(c)

		route, pathParam, err = router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput = &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		err = service.PutRepaidStates(c)
		asr.NoError(err)
		asr.Equal(http.StatusBadRequest, rec.Code)
	})
}
