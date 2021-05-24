package model

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/container/intsets"
	"testing"
	"time"
)

func TestCreateApplicationsDetail(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}

		detail, err := repo.createApplicationsDetail(db, appId, "userId", ApplicationType{Type: Club}, "Title", "Remarks", 1000, time.Now())
		asr.NoError(err)
		asr.Equal(detail.ApplicationID, appId)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = repo.createApplicationsDetail(db, id, "userId", ApplicationType{Type: Club}, "Title", "Remakrs", 1000, time.Now())
		asr.Error(err)
	})
}

func TestPutApplicationsDetail(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, "userId")
		if err != nil {
			panic(err)
		}

		oldDetail, err := repo.createApplicationsDetail(db, appId, "userId", ApplicationType{Type: Club}, "Title", "Remarks", 1000, time.Now())
		if err != nil {
			panic(err)
		}

		newUserId := "user2Id"
		newDetail, err := repo.putApplicationsDetail(db, oldDetail.ID, newUserId, nil, "", "", nil, nil)
		asr.NoError(err)

		asr.Equal(newDetail.ApplicationID, appId)
		asr.Equal(newDetail.UpdateUserTrapID.TrapId, newUserId)
		asr.Equal(newDetail.Type, oldDetail.Type)
		asr.Equal(newDetail.Title, oldDetail.Title)
		asr.Equal(newDetail.Remarks, oldDetail.Remarks)
		asr.Equal(newDetail.Amount, oldDetail.Amount)
		asr.Equal(newDetail.PaidAt.PaidAt, oldDetail.PaidAt.PaidAt.Truncate(time.Hour*24))
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		_, err := repo.putApplicationsDetail(db, intsets.MaxInt, "userId", nil, "", "", nil, nil)
		asr.Error(err)
	})
}