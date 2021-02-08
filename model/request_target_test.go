package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateRequestTarget(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		err = repo.createRequestTarget(db, id, "UserId")
		asr.NoError(err)
	})
}

func TestUpdateRequestTarget(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		err = repo.createRequestTarget(db, id, "UserId")
		if err != nil {
			panic(err)
		}
		dt := time.Now()
		_, _, err = repo.UpdateRequestTarget(id, "UserId", "userId", dt)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}
		dt := time.Now()
		_, _, err = repo.UpdateRequestTarget(id, "UserId", "userId", dt)
		asr.Error(err)
	})
}

func TestDeleteRequestTarget(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		err = repo.createRequestTarget(db, id, "UserId")
		if err != nil {
			panic(err)
		}
		err = repo.createRequestTarget(db, id, "UserId1")
		if err != nil {
			panic(err)
		}

		err = repo.deleteRequestTargetByRequestID(db, id)
		asr.NoError(err)

		err = repo.createRequestTarget(db, id, "UserId")
		asr.NoError(err)
		err = repo.createRequestTarget(db, id, "UserId1")
		asr.NoError(err)
	})
}
