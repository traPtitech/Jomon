package router

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"net/http/httptest"
	"os"
	"testing"
)

var router *openapi3filter.Router

func TestMain(m *testing.M) {
	setRouter()

	code := m.Run()
	os.Exit(code)
}

func setRouter() {
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)
	router = openapi3filter.NewRouter().WithSwaggerFromFile("../docs/swagger.yaml")
}

func validateResponse(ctx *context.Context, requestValidationInput *openapi3filter.RequestValidationInput, rec *httptest.ResponseRecorder) error {
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 rec.Code,
		Header:                 rec.Header(),
	}
	responseValidationInput.SetBodyBytes(rec.Body.Bytes())

	return openapi3filter.ValidateResponse(*ctx, responseValidationInput)
}
