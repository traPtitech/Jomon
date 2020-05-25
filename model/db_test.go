package model

import (
	"math/rand"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

const userNameChrs = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var randSrc = rand.NewSource(0)
var repo = &applicationRepository{}
var commentRepo = &commentRepository{}
var adminRepo = &administratorRepository{}

func TestMain(m *testing.M) {
	db := setupDB()
	code := m.Run()
	deleteAllRecord(db)
	os.Exit(code)
}

func setupDB() *gorm.DB {
	db, err := EstablishConnectionForTest()
	if err != nil {
		panic(err)
	}

	err = Migrate()
	if err != nil {
		panic(err)
	}

	return db
}

func deleteAllRecord(db *gorm.DB) {
	db.Delete(&Administrator{})
	db.Delete(&ApplicationsDetail{})
	db.Delete(&ApplicationsImage{})
	db.Unscoped().Delete(&Comment{})
	db.Delete(&RepayUser{})
	db.Delete(&StatesLog{})
	db.Delete(&Application{})
}

func generateRandomUserName() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = userNameChrs[int(randSrc.Int63()%int64(len(userNameChrs)))]
	}

	return string(b)
}
