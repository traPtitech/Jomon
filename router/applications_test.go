package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/traPtitech/Jomon/model"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"testing"
	"time"
)

const MultipartBoundary = "-------------------------Multipart_Boundary"

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
	sort string,
	currentState *model.StateType,
	financialYear *int,
	applicant string,
	typ *model.ApplicationType,
	submittedSince *time.Time,
	submittedUntil *time.Time,
	giveAdmin bool,
) ([]model.Application, error) {
	ret := m.Called(sort, currentState, financialYear, applicant, typ, submittedSince, submittedUntil, giveAdmin)
	if sort != "" {
		m.asr.Contains([...]string{"", "created_at", "-created_at", "title", "-title"}, sort)
	}

	if applicant != "" {
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
	title string,
	remarks string,
	amount *int,
	paidAt *time.Time,
) error {
	ret := m.Called(appId, updateUserTrapId, typ, title, remarks, amount, paidAt)

	m.asr.NotEqual("", updateUserTrapId)

	if typ != nil {
		m.asr.Contains([...]int{model.Club, model.Contest, model.Event, model.Public}, typ.Type)
	}

	if amount != nil {
		m.asr.Less(0, *amount)
	}

	return ret.Error(0)
}

/*
	Util functions
*/

func GenerateApplication(
	appId uuid.UUID,
	createUserTrapID string,
	typ model.ApplicationType,
	title string,
	remarks string,
	amount int,
	paidAt time.Time,
) model.Application {
	detail := model.ApplicationsDetail{
		ID:            1,
		ApplicationID: appId,
		UpdateUserTrapID: model.User{
			TrapId:  createUserTrapID,
			IsAdmin: false,
		},
		Type:      typ,
		Title:     title,
		Remarks:   remarks,
		Amount:    amount,
		PaidAt:    model.PaidAt{PaidAt: paidAt},
		UpdatedAt: time.Now(),
	}
	state := model.StatesLog{
		ID:            1,
		ApplicationID: appId,
		UpdateUserTrapID: model.User{
			TrapId:  createUserTrapID,
			IsAdmin: false,
		},
		ToState:   model.StateType{Type: model.Submitted},
		Reason:    "",
		CreatedAt: time.Now(),
	}
	return model.Application{
		ID:                       appId,
		LatestApplicationsDetail: detail,
		ApplicationsDetailsID:    1,
		LatestStatesLog:          state,
		LatestStatus:             model.StateType{Type: model.Submitted},
		StatesLogsID:             1,
		CreateUserTrapID: model.User{
			TrapId:  createUserTrapID,
			IsAdmin: false,
		},
		CreatedAt:           time.Now(),
		ApplicationsDetails: []model.ApplicationsDetail{detail},
		StatesLogs:          []model.StatesLog{state},
		ApplicationsImages:  []model.ApplicationsImage{},
		Comments:            []model.Comment{},
		RepayUsers:          []model.RepayUser{},
	}
}

/*
	Function Tests
*/

func TestGetApplication(t *testing.T) {
	appRepMock := NewApplicationRepositoryMock(t)

	appId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	application := GenerateApplication(appId, "User1", model.ApplicationType{Type: model.Club}, "Title", "Remakrs", 10000, time.Now())

	appRepMock.On("GetApplication", application.ID, true, true).Return(application, nil)
	appRepMock.On("GetApplication", mock.Anything, mock.Anything, mock.Anything).Return(model.Application{}, gorm.ErrRecordNotFound)

	service := Service{
		Applications: NewApplicationService(appRepMock),
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/api/applications/"+appId.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId")
		c.SetParamNames("applicationId")
		c.SetParamValues(appId.String())

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

		req := httptest.NewRequest(http.MethodGet, "/api/applications/"+id.String(), nil)
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
	appRepMock := NewApplicationRepositoryMock(t)

	appId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	application := GenerateApplication(appId, "User1", model.ApplicationType{Type: model.Club}, "Title", "Remakrs", 10000, time.Now())

	appRepMock.On("GetApplicationList", "", (*model.StateType)(nil), (*int)(nil), "", (*model.ApplicationType)(nil), (*time.Time)(nil), (*time.Time)(nil), mock.Anything).Return([]model.Application{application}, nil)
	appRepMock.On("GetApplicationList", "title", (*model.StateType)(nil), (*int)(nil), "User1", (*model.ApplicationType)(nil), (*time.Time)(nil), (*time.Time)(nil), true).Return([]model.Application{application}, nil)
	appRepMock.On("GetApplicationList", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.Application{}, nil)

	service := Service{
		Applications: NewApplicationService(appRepMock),
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		t.Run("noQueryParams", func(t *testing.T) {
			asr := assert.New(t)
			e := echo.New()
			ctx := context.TODO()

			req := httptest.NewRequest(http.MethodGet, "/api/applications", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/applications")

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

			err = service.GetApplicationList(c)
			asr.NoError(err)

			asr.Equal(http.StatusOK, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)
		})

		t.Run("giveSortAndApplicant", func(t *testing.T) {
			asr := assert.New(t)
			e := echo.New()
			ctx := context.TODO()

			q := make(url.Values)
			q.Add("sort", "title")
			q.Add("applicant", "User1")
			req := httptest.NewRequest(http.MethodGet, "/api/applications?"+q.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/applications")

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

			err = service.GetApplicationList(c)
			asr.NoError(err)

			asr.Equal(http.StatusOK, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)
		})

		t.Run("noApplicationHit", func(t *testing.T) {
			asr := assert.New(t)
			e := echo.New()
			ctx := context.TODO()

			q := make(url.Values)
			q.Add("sort", "title")
			q.Add("applicant", "User2")
			req := httptest.NewRequest(http.MethodGet, "/api/applications?"+q.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/applications")

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

			err = service.GetApplicationList(c)
			asr.NoError(err)

			asr.Equal(http.StatusOK, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)

			var results []model.Application
			err = json.Unmarshal(rec.Body.Bytes(), &results)
			asr.NoError(err)
			asr.Empty(results)
		})

		t.Run("invalidQueryParameter", func(t *testing.T) {
			asr := assert.New(t)
			e := echo.New()
			ctx := context.TODO()

			q := make(url.Values)
			q.Add("submitted_since", "invalid")
			req := httptest.NewRequest(http.MethodGet, "/api/applications?"+q.Encode(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/applications")

			route, pathParam, err := router.FindRoute(req.Method, req.URL)
			if err != nil {
				// panic(err)
			}

			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    req,
				PathParams: pathParam,
				Route:      route,
			}

			if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
				// panic(err)
			}

			err = service.GetApplicationList(c)
			asr.NoError(err)

			asr.Equal(http.StatusBadRequest, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)
		})
	})
}

func TestPostApplication(t *testing.T) {
	appRepMock := NewApplicationRepositoryMock(t)

	title := "夏コミの交通費をお願いします。"
	remarks := "〇〇駅から〇〇駅への移動"
	amount := 1000
	paidAt := time.Now().Round(time.Second)

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	appRepMock.On("GetApplication", id, mock.Anything, mock.Anything).Return(GenerateApplication(id, "UserId", model.ApplicationType{Type: model.Club}, title, remarks, amount, paidAt), nil)
	appRepMock.On("BuildApplication", "UserId", model.ApplicationType{Type: model.Club}, title, remarks, amount, mock.Anything).Return(id, nil)

	service := Service{
		Applications: NewApplicationService(appRepMock),
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := &bytes.Buffer{}
		mpw := multipart.NewWriter(body)
		if err := mpw.SetBoundary(MultipartBoundary); err != nil {
			panic(err)
		}

		part := make(textproto.MIMEHeader)
		part.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, "details"))
		part.Set("Content-Type", "application/json")
		writer, err := mpw.CreatePart(part)
		if err != nil {
			panic(err)
		}
		_, err = writer.Write([]byte(fmt.Sprintf(`
			{
				"type": "club",
				"title": "%s",
				"remarks": "%s",
				"paid_at": "%s",
				"amount": %d,
				"repaid_to_id": [
					"User1"
				]
			}
		`, title, remarks, paidAt.Format(time.RFC3339), amount)))
		if err != nil {
			panic(err)
		}

		if err = mpw.Close(); err != nil {
			panic(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/applications", body)
		req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", MultipartBoundary))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications")

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

		err = service.PostApplication(c)
		asr.NoError(err)

		asr.Equal(http.StatusCreated, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := &bytes.Buffer{}
		mpw := multipart.NewWriter(body)
		if err := mpw.SetBoundary(MultipartBoundary); err != nil {
			panic(err)
		}

		part := make(textproto.MIMEHeader)
		part.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, "details"))
		part.Set("Content-Type", "application/json")
		writer, err := mpw.CreatePart(part)
		if err != nil {
			panic(err)
		}
		_, err = writer.Write([]byte(fmt.Sprintf(`
			{
				"type": "club",
				"title": "%s",
				"remarks": "%s",
				"paid_at": "%s",
				"amount": %d
			}
		`, title, remarks, paidAt.Format(time.RFC3339), amount)))
		if err != nil {
			// panic(err)
		}

		if err = mpw.Close(); err != nil {
			panic(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/applications", body)
		req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", MultipartBoundary))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications")

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			// panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			// panic(err)
		}

		err = service.PostApplication(c)
		asr.NoError(err)

		asr.Equal(http.StatusBadRequest, rec.Code)
	})
}

func TestPatchApplication(t *testing.T) {

}
