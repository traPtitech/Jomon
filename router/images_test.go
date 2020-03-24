package router

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/traPtitech/Jomon/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type applicationsImageRepositoryMock struct {
	mock.Mock
}

func (m *applicationsImageRepositoryMock) CreateApplicationsImage(applicationId uuid.UUID, src io.Reader, mimeType string) (model.ApplicationsImage, error) {
	ret := m.Called(applicationId, src, mimeType)
	return ret.Get(0).(model.ApplicationsImage), ret.Error(1)
}

func (m *applicationsImageRepositoryMock) GetApplicationsImage(id uuid.UUID) (model.ApplicationsImage, error) {
	ret := m.Called(id)
	return ret.Get(0).(model.ApplicationsImage), ret.Error(1)
}

func (m *applicationsImageRepositoryMock) OpenApplicationsImage(appImg model.ApplicationsImage) (io.ReadCloser, error) {
	ret := m.Called(appImg)
	return ret.Get(0).(io.ReadCloser), ret.Error(1)
}

func (m *applicationsImageRepositoryMock) DeleteApplicationsImage(appImg model.ApplicationsImage) error {
	ret := m.Called(appImg)
	return ret.Error(0)
}

func TestGetImages(t *testing.T) {
	t.Parallel()

	imRepMock := new(applicationsImageRepositoryMock)

	imId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	appId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	sampleIm := model.ApplicationsImage{
		ID:            imId,
		ApplicationID: appId,
		MimeType:      "image/png",
		CreatedAt:     time.Now(),
	}

	errImId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	imRepMock.On("GetApplicationsImage", imId).Return(sampleIm, nil)
	imRepMock.On("GetApplicationsImage", errImId).Return(model.ApplicationsImage{}, gorm.ErrRecordNotFound)

	userRepMock := NewUserRepositoryMock(t, "UserId", "AdminUserId")

	service := Service{
		Administrators: nil,
		Applications:   nil,
		Comments:       nil,
		Images:         imRepMock,
		Users:          userRepMock,
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := newEcho(service)
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/api/images/"+imId.String(), nil)
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("/images/:imageId")
		c.SetParamNames("imageId")
		c.SetParamValues(imId.String())

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

		err = service.GetImages(c)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := newEcho(service)
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodGet, "/api/images/"+errImId.String(), nil)
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("/images/:imageId")
		c.SetParamNames("imageId")
		c.SetParamValues(errImId.String())

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

		err = service.GetImages(c)
		asr.NoError(err)
		asr.Equal(http.StatusNotFound, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})
}

func TestDeleteImages(t *testing.T) {
	t.Parallel()

	appRepMock := NewApplicationRepositoryMock(t)

	appId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	userId := "UserId"

	application := GenerateApplication(appId, userId, model.ApplicationType{Type: model.Club}, "Title", "Remarks", 10000, time.Now())

	appRepMock.On("GetApplication", application.ID, true).Return(application, nil)

	adminRepMock := NewAdministratorRepositoryMock("AdminUserId")

	userRepMock := NewUserRepositoryMock(t, userId, "AdminUserId")

	anotherToken := "AnotherToken"
	userRepMock.On("GetMyUser", anotherToken).Return(model.User{TrapId: "AnotherId"}, nil)

	imRepMock := new(applicationsImageRepositoryMock)

	imId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	sampleIm := model.ApplicationsImage{
		ID:            imId,
		ApplicationID: appId,
		MimeType:      "image/png",
		CreatedAt:     time.Now(),
	}

	errImId, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	imRepMock.On("GetApplicationsImage", imId).Return(sampleIm, nil)
	imRepMock.On("GetApplicationsImage", errImId).Return(model.ApplicationsImage{}, gorm.ErrRecordNotFound)

	imRepMock.On("DeleteApplicationsImage", sampleIm).Return(nil)

	service := Service{
		Administrators: adminRepMock,
		Applications:   appRepMock,
		Comments:       nil,
		Images:         imRepMock,
		Users:          userRepMock,
	}

	t.Run("shouldSuccess", func(t *testing.T) {
		t.Parallel()

		t.Run("deleteByApplicationAuthor", func(t *testing.T) {
			asr := assert.New(t)
			e := newEcho(service)
			ctx := context.TODO()

			req := httptest.NewRequest(http.MethodDelete, "/api/images/"+imId.String(), nil)
			req.Header.Set("Authorization", userRepMock.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("/images/:imageId")
			c.SetParamNames("imageId")
			c.SetParamValues(imId.String())

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

			err = service.DeleteImages(c)
			asr.NoError(err)
			asr.Equal(http.StatusNoContent, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)
		})

		t.Run("deleteByAdmin", func(t *testing.T) {
			asr := assert.New(t)
			e := newEcho(service)
			ctx := context.TODO()

			req := httptest.NewRequest(http.MethodDelete, "/api/images/"+imId.String(), nil)
			req.Header.Set("Authorization", userRepMock.adminToken)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("/images/:imageId")
			c.SetParamNames("imageId")
			c.SetParamValues(imId.String())

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

			err = service.DeleteImages(c)
			asr.NoError(err)
			asr.Equal(http.StatusNoContent, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)
		})
	})

	t.Run("shouldFail", func(t *testing.T) {
		t.Parallel()

		t.Run("notExistImage", func(t *testing.T) {
			asr := assert.New(t)
			e := newEcho(service)
			ctx := context.TODO()

			req := httptest.NewRequest(http.MethodDelete, "/api/images/"+errImId.String(), nil)
			req.Header.Set("Authorization", userRepMock.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("/images/:imageId")
			c.SetParamNames("imageId")
			c.SetParamValues(errImId.String())

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

			err = service.DeleteImages(c)
			asr.NoError(err)
			asr.Equal(http.StatusNotFound, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)
		})

		t.Run("forbidden", func(t *testing.T) {
			asr := assert.New(t)
			e := newEcho(service)
			ctx := context.TODO()

			req := httptest.NewRequest(http.MethodDelete, "/api/images/"+imId.String(), nil)
			req.Header.Set("Authorization", anotherToken)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("/images/:imageId")
			c.SetParamNames("imageId")
			c.SetParamValues(errImId.String())

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

			err = service.DeleteImages(c)
			asr.NoError(err)
			asr.Equal(http.StatusForbidden, rec.Code)

			err = validateResponse(&ctx, requestValidationInput, rec)
			asr.NoError(err)
		})
	})
}
