package model

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type transactionRepositoryMock struct {
	mock.Mock
	token string
}

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		transactionID, err := transactionRepo.CreateTransaction(100, []string{"target1"}, &id)
		asr.NoError(err)
		asr.NotEqual(transactionID, uuid.Nil)
	})

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		transactionID, err := transactionRepo.CreateTransaction(100, []string{"target1"}, nil)
		asr.NoError(err)
		asr.NotEqual(transactionID, uuid.Nil)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = transactionRepo.CreateTransaction(100, []string{"target1"}, &id)
		asr.Error(err)
	})
}

func TestGetTransaction(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		trnsID, err := transactionRepo.CreateTransaction(100, []string{"target1"}, nil)
		if err != nil {
			panic(err)
		}

		getTrns, err := transactionRepo.GetTransaction(trnsID)

		asr.NoError(err)
		asr.Equal(getTrns.ID, trnsID)
		detail, err := transactionRepo.GetTransactionDetail(getTrns.ID)
		if err != nil {
			panic(err)
		}
		asr.Equal(detail.Amount, 100)
		asr.Equal(detail.Target, "target1")
	})
}

func TestPatchTransaction(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		amount := 1000
		amount2 := 10000
		target := "target1"
		// TODO targetを配列に変更
		trnsID, err := transactionRepo.CreateTransaction(amount, []string{target}, nil)
		if err != nil {
			panic(err)
		}

		err = transactionRepo.PatchTransaction(trnsID, &amount2, []string{""}, nil)
		asr.NoError(err)

		trns, err := transactionRepo.GetTransaction(trnsID)
		if err != nil {
			panic(err)
		}

		detail, err := transactionRepo.GetTransactionDetail(trns.ID)
		if err != nil {
			panic(err)
		}
		asr.Equal(amount2, detail.Amount)
		asr.Equal(target, detail.Target)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		err = transactionRepo.PatchTransaction(id, nil, []string{generateRandomUserName()}, nil)
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})

}

func TestGetTransactionList(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		deleteAllRecord(db)

		target1 := "Target1"
		target2 := "Target2"
		target3 := "Target3"

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		app1SubTime := time.Date(2020, 4, 10, 12, 0, 0, 0, time.Local)
		app1Id := createTransactionWithCreatedTime(100, target1, id, app1SubTime)

		app2SubTime := time.Date(2020, 4, 20, 12, 0, 0, 0, time.Local)
		app2Id := createTransactionWithCreatedTime(100, target2, id, app2SubTime)

		app3SubTime := time.Date(2020, 4, 30, 12, 0, 0, 0, time.Local)
		app3Id := createTransactionWithCreatedTime(100, target3, id, app3SubTime)
		app3, err := transactionRepo.GetTransaction(app3Id)
		if err != nil {
			panic(err)
		}

		app4SubTime := time.Date(2019, 4, 10, 12, 0, 0, 0, time.Local)
		app4Id := createTransactionWithCreatedTime(100, target1, id, app4SubTime)

		amount2 := 1000
		_ = transactionRepo.PatchTransaction(app3.ID, &amount2, []string{}, nil)
		if err != nil {
			panic(err)
		}

		t.Parallel()

		t.Run("allNil", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := transactionRepo.GetTransactionList("", nil, "", nil, nil)
			asr.NoError(err)

			asr.Len(apps, 4)
			asr.Equal([]uuid.UUID{app3Id, app2Id, app1Id, app4Id}, mapToTransactionID(apps))

			asr.False(apps[0].CreatedAt.Before(apps[1].CreatedAt))
			asr.False(apps[1].CreatedAt.Before(apps[2].CreatedAt))

			detail, err := transactionRepo.GetTransactionDetail(apps[2].ID)
			asr.NoError(err)
			asr.Equal(target3, detail.Target)
		})

		t.Run("filterByFinancialYear", func(t *testing.T) {
			asr := assert.New(t)
			financialYear := 2019

			apps, err := transactionRepo.GetTransactionList("", &financialYear, "", nil, nil)
			asr.NoError(err)

			asr.Len(apps, 1)
			asr.Equal(app4Id, apps[0].ID)
		})

		t.Run("filterByApplicant", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := transactionRepo.GetTransactionList("", nil, target2, nil, nil)
			asr.NoError(err)

			asr.Len(apps, 2)
		})

		t.Run("emptyResult", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := transactionRepo.GetTransactionList("", nil, target3, nil, nil)
			asr.NoError(err)

			asr.Empty(apps)
		})

		t.Run("filterBySubmitted", func(t *testing.T) {
			beforeApp2 := app2SubTime.Add(-24 * time.Hour)
			beforeApp3 := app3SubTime.Add(-24 * time.Hour)

			t.Parallel()

			t.Run("since", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := transactionRepo.GetTransactionList("", nil, "", nil, &beforeApp3)
				asr.NoError(err)

				asr.Len(apps, 1)
				asr.Equal(app3Id, apps[0].ID)
			})

			t.Run("until", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := transactionRepo.GetTransactionList("", nil, "", nil, &beforeApp3)
				asr.NoError(err)

				asr.Len(apps, 3)
				asr.Equal([]uuid.UUID{app2Id, app1Id, app4Id}, mapToTransactionID(apps))
			})

			t.Run("both", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := transactionRepo.GetTransactionList("", nil, "", &beforeApp2, &beforeApp3)
				asr.NoError(err)

				asr.Len(apps, 1)
				asr.Equal(app2Id, apps[0].ID)
			})
		})

		t.Run("sort", func(t *testing.T) {
			tests := []struct {
				SortBy string
				Should []uuid.UUID
			}{
				{
					SortBy: "created_at",
					Should: []uuid.UUID{app3Id, app2Id, app1Id, app4Id},
				},
				{
					SortBy: "-created_at",
					Should: []uuid.UUID{app4Id, app1Id, app2Id, app3Id},
				},
				{
					SortBy: "title",
					Should: []uuid.UUID{app2Id, app3Id, app1Id, app4Id},
				},
				{
					SortBy: "-title",
					Should: []uuid.UUID{app4Id, app1Id, app3Id, app2Id},
				},
			}

			t.Parallel()

			for _, test := range tests {
				t.Run(test.SortBy, func(t *testing.T) {
					asr := assert.New(t)

					apps, err := transactionRepo.GetTransactionList(test.SortBy, nil, "", nil, nil)
					asr.NoError(err)

					asr.Len(apps, 4)
					asr.Equal(test.Should, mapToTransactionID(apps))
				})
			}
		})
	})
}

func mapToTransactionID(trnss []Transaction) []uuid.UUID {
	trnsIds := make([]uuid.UUID, len(trnss))
	for i := range trnss {
		trnsIds[i] = trnss[i].ID
	}

	return trnsIds
}

func createTransactionWithCreatedTime(amount int, target string, requestID uuid.UUID, createdAt time.Time) uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	err = db.Create(&Transaction{
		ID:        id,
		CreatedAt: createdAt,
	}).Error
	if err != nil {
		panic(err)
	}

	detailID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	err = db.Create(&TransactionDetail{
		ID:            detailID,
		TransactionID: id,
		Amount:        amount,
		Target:        target,
		RequestID:     requestID,
	}).Error
	if err != nil {
		panic(err)
	}

	return id
}
