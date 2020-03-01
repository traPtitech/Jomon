package router

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/traPtitech/Jomon/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

/*
	Definitions of mock
*/

type applicationRepositoryMock struct {
	mock.Mock
	asr *assert.Assertions
}

func NewApplicationRepositoryMock(t *testing.T) *applicationRepositoryMock {
	m := new(applicationRepositoryMock)
	m.asr = assert.New(t)
	return m
}

func (m *applicationRepositoryMock) GetApplication(id uuid.UUID, giveAdmin bool, preload bool) (model.Application, error) {
	ret := m.Called(id, giveAdmin, preload)
	return ret.Get(0).(model.Application), ret.Error(1)
}

func (m *applicationRepositoryMock) GetApplicationList(
	sort *string,
	currentState *model.StateType,
	financialYear *int,
	applicant *string,
	typ *model.ApplicationType,
	submittedSince *time.Time,
	submittedUntil *time.Time,
	giveAdmin bool,
) ([]model.Application, error) {
	ret := m.Called(sort, currentState, financialYear, applicant, typ, submittedSince, submittedUntil, giveAdmin)
	if sort != nil {
		m.asr.Contains([...]string{"", "created_at", "-created_at", "title", "-title"}, *sort)
	}

	if applicant != nil {
		m.asr.NotEqual("", applicant)
	}

	if typ != nil {
		m.asr.Contains([...]int{model.Club, model.Contest, model.Event, model.Public}, typ.Type)
	}

	return ret.Get(0).([]model.Application), ret.Error(1)
}

func (m *applicationRepositoryMock) BuildApplication(
	createUserTrapID string,
	typ model.ApplicationType,
	title string,
	remarks string,
	amount int,
	paidAt time.Time,
) (uuid.UUID, error) {
	ret := m.Called(createUserTrapID, typ, title, remarks, amount, paidAt)
	m.asr.NotEqual("", createUserTrapID)
	m.asr.Contains([...]int{model.Club, model.Contest, model.Event, model.Public}, typ.Type)
	m.asr.NotEqual("", title)
	m.asr.Less(0, amount)
	return ret.Get(0).(uuid.UUID), ret.Error(1)
}

func (m *applicationRepositoryMock) PatchApplication(
	appId uuid.UUID,
	updateUserTrapId string,
	typ *model.ApplicationType,
	title *string,
	remarks *string,
	amount *int,
	paidAt *time.Time,
) error {
	ret := m.Called(appId, updateUserTrapId, typ, title, remarks, amount, paidAt)

	m.asr.NotEqual("", updateUserTrapId)

	if typ != nil {
		m.asr.Contains([...]int{model.Club, model.Contest, model.Event, model.Public}, typ.Type)
	}

	if title != nil {
		m.asr.NotEqual("", *title)
	}

	if amount != nil {
		m.asr.Less(0, *amount)
	}

	return ret.Error(0)
}

/*
	Function Tests
*/

func TestGetApplication(t *testing.T) {
	appId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	appRepMock := NewApplicationRepositoryMock(t)

	detail := model.ApplicationsDetail{
		ID:            1,
		ApplicationID: appId,
		UpdateUserTrapID: model.User{
			TrapId:  "User1",
			IsAdmin: false,
		},
		Type:      model.ApplicationType{Type: model.Contest},
		Title:     "title",
		Remarks:   "remarks",
		Amount:    10000,
		PaidAt:    model.PaidAt{PaidAt: time.Now()},
		UpdatedAt: time.Now(),
	}
	state := model.StatesLog{
		ID:            1,
		ApplicationID: appId,
		UpdateUserTrapID: model.User{
			TrapId:  "User1",
			IsAdmin: false,
		},
		ToState:   model.StateType{Type: model.Submitted},
		Reason:    "",
		CreatedAt: time.Now(),
	}
	application := model.Application{
		ID:                       appId,
		LatestApplicationsDetail: detail,
		ApplicationsDetailsID:    1,
		LatestStatesLog:          state,
		LatestStatus:             model.StateType{},
		StatesLogsID:             0,
		CreateUserTrapID:         model.User{},
		CreatedAt:                time.Time{},
		ApplicationsDetails:      []model.ApplicationsDetail{detail},
		StatesLogs:               []model.StatesLog{state},
		ApplicationsImages:       []model.ApplicationsImage{},
		Comments:                 []model.Comment{},
		RepayUsers: []model.RepayUser{{
			ID:            1,
			ApplicationID: appId,
			RepaidToUserTrapID: model.User{
				TrapId:  "User1",
				IsAdmin: false,
			},
			RepaidByUserTrapID: nil,
			RepaidAt:           nil,
		}},
	}

	appRepMock.On("GetApplication", appId, true, true).Return(application, nil)
	appRepMock.On("GetApplication", mock.Anything, mock.Anything, mock.Anything).Return(model.Application{}, gorm.ErrRecordNotFound)

	service := Service{
		Applications: NewApplicationService(NewApplicationRepositoryMock(t)),
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId")
		c.SetParamNames("applicationId")
		c.SetParamValues(appId.String())

		route, pathParam, err := router.FindRoute(req.Method, &url.URL{Path: "/api/applications/" + appId.String()})
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

		err = service.GetApplication(c)
		asr.NoError(err)

		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("/applications/:applicationId")
		c.SetParamNames("applicationId")
		c.SetParamValues(id.String())

		route, pathParam, err := router.FindRoute(req.Method, &url.URL{Path: "/api/applications/" + id.String()})
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

		err = service.GetApplication(c)
		asr.NoError(err)

		asr.Equal(http.StatusNotFound, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
}

func TestGetApplicationList(t *testing.T) {
	//t.Run("shouldSuccess", func(t *testing.T) {
	//	asr := assert.New(t)
	//	e := echo.New()
	//	ctx := context.TODO()
	//
	//
	//})
}

func TestPostApplication(t *testing.T) {

}

func TestPatchApplication(t *testing.T) {

}
