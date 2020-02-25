package router

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/traPtitech/Jomon/model"
	"os"
	"testing"
)

var router *openapi3filter.Router

func TestMain(m *testing.M) {
	if _, err := model.EstablishConnection(); err != nil {
		panic(err)
	}
	if err := model.Migrate(); err != nil {
		panic(err)
	}

	setRouter()

	code := m.Run()
	os.Exit(code)
}

func setRouter() {
	openapi3.DefineStringFormat("uuid", openapi3.FormatOfStringForUUIDOfRFC4122)
	router = openapi3filter.NewRouter().WithSwaggerFromFile("../docs/swagger.yaml")
}
