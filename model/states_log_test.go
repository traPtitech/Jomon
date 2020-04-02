package model

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateStatesLog(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}

		state, err := repo.createStatesLog(db, id, "userId")
		asr.NoError(err)
		asr.Equal(state.ApplicationID, id)
		asr.Equal(state.ToState.Type, Submitted)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = repo.createStatesLog(db, id, "userId")
		asr.Error(err)
	})
}

func TetsUpdateStatesLog(t *testing.T){
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)


		id, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}

		_, err = repo.createStatesLog(db, id, "userId")
		if err != nil {
			panic(err)
		}

		state, err := repo.UpdateStatesLog(id, "userId", "reason", StateType{Type: Accepted})
		asr.NoError(err)
		asr.Equal(state.ApplicationID, id)
		asr.Equal(state.ToState.Type, Accepted)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		
		id, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}


		_, err = repo.UpdateStatesLog(id, "userId", "reason", StateType{Type: Accepted})
		asr.Error(err)
	})
}