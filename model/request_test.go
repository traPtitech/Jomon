package model

import (
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRequest(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		id, err := repo.createRequest(db, "userId")
		if err != nil {
			panic(err)
		}

		state, err := repo.createRequestStatus(db, id, "userId")
		asr.NoError(err)
		asr.Equal(state.RequestID, id)
		asr.Equal(state.Status, Submitted)
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

func TestBuildRequest(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appID, err := repo.BuildRequest("User1", "Title", "Remarks", 1000, time.Now(), []string{"User1"})
		asr.NoError(err)
		asr.NotEqual(appID, uuid.Nil)
	})
}

func TestGetRequest(t *testing.T) {
	t.Parallel()

	sm := new(storageMock)
	sm.On("Save", mock.Anything, mock.Anything).Return(nil)

	fileRepo := NewFileRepository(sm)

	t.Run("shouldSuccess?giveAdmin=true&preload=true", func(t *testing.T) {
		asr := assert.New(t)

		user := generateRandomUserName()

		appID, err := repo.BuildRequest(user, "Title", "Remarks", 1000, time.Now(), []string{user})
		if err != nil {
			panic(err)
		}

		comment, err := commentRepo.CreateComment(appID, "This is comment.", user)

		img, err := fileRepo.CreateFile(appID, strings.NewReader("TestData"), "image/png")

		app, err := repo.GetRequest(appID, true)

		asr.NoError(err)
		asr.Equal(appID, app.ID)

		asr.Equal(app.RequestStatusID, app.LatestRequestStatus.ID)
		asr.Len(app.RequestStatus, 1)
		asr.Equal(app.LatestRequestStatus, app.RequestStatus[0])
		asr.Len(app.RequestTargets, 1)
		asr.Len(app.Files, 1)

		asr.Equal(img.ID, app.Files[0].ID)
		asr.Equal(img.MimeType, app.Files[0].MimeType)
		asr.WithinDuration(img.CreatedAt, app.Files[0].CreatedAt, 1*time.Second)

		asr.Equal(comment.ID, app.Comments[0].ID)
		asr.Equal(comment.Comment, app.Comments[0].Comment)
		asr.Equal(comment.UserTrapID, app.Comments[0].UserTrapID)
		asr.Len(app.Comments, 1)
	})

	t.Run("shouldSuccess?giveAdmin=true&preload=false", func(t *testing.T) {
		asr := assert.New(t)

		user := generateRandomUserName()

		appID, err := repo.BuildRequest(user, "Title", "Remarks", 1000, time.Now(), []string{user})
		if err != nil {
			panic(err)
		}

		app, err := repo.GetRequest(appID, false)

		asr.NoError(err)
		asr.Equal(appID, app.ID)

	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = repo.GetRequest(id, true)
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})
}

func TestPatchRequest(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		title := "Title"
		remarks := "Remarks"
		amount := 1000
		paidAt := time.Now()
		appID, err := repo.BuildRequest("User", title, remarks, amount, paidAt, []string{"User"})
		if err != nil {
			panic(err)
		}

		err = repo.PatchRequest(appID, "User", "", "", nil, nil, []string{"OtherUser"})
		asr.NoError(err)

		app, err := repo.GetRequest(appID, true)
		if err != nil {
			panic(err)
		}

		asr.Equal(title, app.Title)
		asr.Equal(remarks, app.Content)
		asr.Equal(amount, app.Amount)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		err = repo.PatchRequest(id, generateRandomUserName(), "", "", nil, nil, []string{})
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})

}

func TestGetRequestList(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		deleteAllRecord(db)

		user1 := "User1"
		user2 := "User2"
		user3 := "User3"

		app1SubTime := time.Date(2020, 4, 10, 12, 0, 0, 0, time.Local)
		app1Id := buildRequestWithSubmitTime(user1, app1SubTime, "CCCCC", "Remarks", 10000, time.Now())

		app2SubTime := time.Date(2020, 4, 20, 12, 0, 0, 0, time.Local)
		app2Id := buildRequestWithSubmitTime(user2, app2SubTime, "AAAAA", "Remarks", 10000, time.Now())

		app3SubTime := time.Date(2020, 4, 30, 12, 0, 0, 0, time.Local)
		app3Id := buildRequestWithSubmitTime(user2, app3SubTime, "BBBBB", "Remarks", 10000, time.Now())
		app3, err := repo.GetRequest(app3Id, true)
		if err != nil {
			panic(err)
		}

		app4SubTime := time.Date(2019, 4, 10, 12, 0, 0, 0, time.Local)
		app4Id := buildRequestWithSubmitTime(user1, app4SubTime, "DDDDD", "Remarks", 10000, time.Now())

		// TODO Use a appropriate function defined in model/states_log.go after implementing such a function.
		db.Model(&app3.LatestRequestStatus).Updates(RequestStatus{
			Status: FullyRepaid,
		})

		t.Parallel()

		t.Run("allNil", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetRequestList("", nil, nil, "", nil, nil)
			asr.NoError(err)

			asr.Len(apps, 4)
			asr.Equal([]uuid.UUID{app3Id, app2Id, app1Id, app4Id}, mapToRequestID(apps))

			for _, app := range apps {
				asr.NotZero(app.LatestRequestStatus)
				asr.NotZero(app.LatestStatus)
			}

			asr.False(apps[0].CreatedAt.Before(apps[1].CreatedAt))
			asr.False(apps[1].CreatedAt.Before(apps[2].CreatedAt))

			asr.Equal(user1, apps[2].CreatedBy.TrapID)
		})

		t.Run("filterByFinancialYear", func(t *testing.T) {
			asr := assert.New(t)
			financialYear := 2019

			apps, err := repo.GetRequestList("", nil, &financialYear, "", nil, nil)
			asr.NoError(err)

			asr.Len(apps, 1)
			asr.Equal(app4Id, apps[0].ID)
		})

		t.Run("filterByApplicant", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetRequestList("", nil, nil, user2, nil, nil)
			asr.NoError(err)

			asr.Len(apps, 2)
		})

		t.Run("filterByCurrentState", func(t *testing.T) {
			asr := assert.New(t)

			fr := FullyRepaid
			apps, err := repo.GetRequestList("", &fr, nil, "", nil, nil)
			asr.NoError(err)

			asr.Len(apps, 1)
			asr.Equal(app3Id, apps[0].ID)
		})

		t.Run("emptyResult", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetRequestList("", nil, nil, user3, nil, nil)
			asr.NoError(err)

			asr.Empty(apps)
		})

		t.Run("filterBySubmitted", func(t *testing.T) {
			beforeApp2 := app2SubTime.Add(-24 * time.Hour)
			beforeApp3 := app3SubTime.Add(-24 * time.Hour)

			t.Parallel()

			t.Run("Since", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := repo.GetRequestList("", nil, nil, "", nil, &beforeApp3)
				asr.NoError(err)

				asr.Len(apps, 1)
				asr.Equal(app3Id, apps[0].ID)
			})

			t.Run("until", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := repo.GetRequestList("", nil, nil, "", nil, &beforeApp3)
				asr.NoError(err)

				asr.Len(apps, 3)
				asr.Equal([]uuid.UUID{app2Id, app1Id, app4Id}, mapToRequestID(apps))
			})

			t.Run("both", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := repo.GetRequestList("", nil, nil, "", &beforeApp2, &beforeApp3)
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

					apps, err := repo.GetRequestList(test.SortBy, nil, nil, "", nil, nil)
					asr.NoError(err)

					asr.Len(apps, 4)
					asr.Equal(test.Should, mapToRequestID(apps))
				})
			}
		})
	})
}

func mapToRequestID(apps []Request) []uuid.UUID {
	appIds := make([]uuid.UUID, len(apps))
	for i := range apps {
		appIds[i] = apps[i].ID
	}

	return appIds
}

func buildRequestWithSubmitTime(createUserTrapID string, submittedAt time.Time, title string, remarks string, amount int, paidAt time.Time) uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	err = db.Create(&Request{
		ID:        id,
		CreatedBy: TrapUser{TrapID: createUserTrapID},
		CreatedAt: submittedAt,
	}).Error
	if err != nil {
		panic(err)
	}

	state, err := repo.createRequestStatus(db, id, createUserTrapID)
	if err != nil {
		panic(err)
	}

	err = db.Model(Request{}).Where(&Request{ID: id}).Updates(Request{
		RequestStatusID: state.ID,
	}).Error

	if err != nil {
		panic(err)
	}

	return id
}
