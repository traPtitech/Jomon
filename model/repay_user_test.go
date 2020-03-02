package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
