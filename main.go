package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	e := echo.New()

	e.GET("/", genRootHandler(err == nil))

	e.Start(":1323")
}

func genRootHandler(b bool) func(echo.Context) error {
	return func(c echo.Context) error {
		if b {
			return c.String(http.StatusOK, "Succeeded in access db.")
		} else {
			return c.String(http.StatusOK, "Failed to access db.")
		}
	}
}
