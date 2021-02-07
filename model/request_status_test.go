package model

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequestStatus(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		status, err := repo.createRequestStatus(db, id, "userId")
		asr.NoError(err)
		asr.Equal(status.RequestID, id)
		asr.Equal(status.Status, Submitted)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = repo.createRequestStatus(db, id, "userId")
		asr.Error(err)
	})
}

func TestUpdateRequestStatus(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		_, err = repo.createRequestStatus(db, id, "userId")
		if err != nil {
			panic(err)
		}

		status, err := repo.UpdateRequestStatus(id, "userId", "reason", Accepted)
		asr.NoError(err)
		asr.Equal(status.RequestID, id)
		asr.Equal(status.Status, Accepted)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = repo.UpdateRequestStatus(id, "userId", "reason", Accepted)
		asr.Error(err)
	})
}
