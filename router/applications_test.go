package router

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGetApplication(t *testing.T) {
	appId, err := model.BuildApplication("User", model.ApplicationType{Type: 1}, "Title", "Remarks", 1000, time.Now())
	if err != nil {
		panic(err)
	}

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

		err = GetApplication(c)
		asr.NoError(err)

		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: requestValidationInput,
			Status:                 rec.Code,
			Header:                 rec.Header(),
		}

		responseValidationInput.SetBodyBytes(rec.Body.Bytes())

		err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
		asr.NoError(err)
	})

}
