package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepayUser(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}

		err = repo.createRepayUser(db, id, "UserId")
		asr.NoError(err)
	})
}

func TestUpdateRepayUser(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}

		err = repo.createRepayUser(db, id, "UserId")
		if err != nil {
			panic(err)
		}
		dt := time.Now()
		_, _, err = repo.UpdateRepayUser(id, "UserId", "userId", dt)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}
		dt := time.Now()
		_, _, err = repo.UpdateRepayUser(id, "UserId", "userId", dt)
		asr.Error(err)
	})
}

func TestDeleteRepayUser(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}

		err = repo.createRepayUser(db, id, "UserId")
		if err != nil {
			panic(err)
		}
		err = repo.createRepayUser(db, id, "UserId1")
		if err != nil {
			panic(err)
		}

		err = repo.deleteRepayUserByApplicationID(db, id)
		asr.NoError(err)

		err = repo.createRepayUser(db, id, "UserId")
		asr.NoError(err)
		err = repo.createRepayUser(db, id, "UserId1")
		asr.NoError(err)
	})
}
