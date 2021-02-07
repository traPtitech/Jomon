package model

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateApplication(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, "User1")
		asr.NotEqual(appId, uuid.Nil)
		asr.NoError(err)
	})
}

func TestBuildApplication(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.BuildApplication("User1", ApplicationType{Type: Contest}, "Title", "Remarks", 1000, time.Now(), []string{"User1"})
		asr.NoError(err)
		asr.NotEqual(appId, uuid.Nil)
	})
}

func TestGetApplication(t *testing.T) {
	t.Parallel()

	sm := new(storageMock)
	sm.On("Save", mock.Anything, mock.Anything).Return(nil)

	imageRepo := NewApplicationsImageRepository(sm)

	t.Run("shouldSuccess?giveAdmin=true&preload=true", func(t *testing.T) {
		asr := assert.New(t)

		user := generateRandomUserName()

		appId, err := repo.BuildApplication(user, ApplicationType{Type: Contest}, "Title", "Remarks", 1000, time.Now(), []string{user})
		if err != nil {
			panic(err)
		}

		comment, err := commentRepo.CreateComment(appId, "This is comment.", user)

		img, err := imageRepo.CreateApplicationsImage(appId, strings.NewReader("TestData"), "image/png")

		app, err := repo.GetApplication(appId, true)

		asr.NoError(err)
		asr.Equal(appId, app.ID)

		asr.Equal(app.ApplicationsDetailsID, app.LatestApplicationsDetail.ID)
		asr.Equal(app.StatesLogsID, app.LatestStatesLog.ID)
		asr.Len(app.ApplicationsDetails, 1)
		asr.Equal(app.LatestApplicationsDetail, app.ApplicationsDetails[0])
		asr.Len(app.StatesLogs, 1)
		asr.Equal(app.LatestStatesLog, app.StatesLogs[0])
		asr.Len(app.RepayUsers, 1)
		asr.Len(app.ApplicationsImages, 1)

		asr.Equal(img.ID, app.ApplicationsImages[0].ID)
		asr.Equal(img.MimeType, app.ApplicationsImages[0].MimeType)
		asr.WithinDuration(img.CreatedAt, app.ApplicationsImages[0].CreatedAt, 1*time.Second)

		asr.Equal(comment.ID, app.Comments[0].ID)
		asr.Equal(comment.Comment, app.Comments[0].Comment)
		asr.Equal(comment.UserTrapID, app.Comments[0].UserTrapID)
		asr.Len(app.Comments, 1)
	})

	t.Run("shouldSuccess?giveAdmin=true&preload=false", func(t *testing.T) {
		asr := assert.New(t)

		user := generateRandomUserName()

		appId, err := repo.BuildApplication(user, ApplicationType{Type: Contest}, "Title", "Remarks", 1000, time.Now(), []string{user})
		if err != nil {
			panic(err)
		}

		app, err := repo.GetApplication(appId, false)

		asr.NoError(err)
		asr.Equal(appId, app.ID)

	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = repo.GetApplication(id, true)
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})
}

func TestPatchApplication(t *testing.T) {
	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		typ := ApplicationType{Type: Club}
		title := "Title"
		remarks := "Remarks"
		amount := 1000
		paidAt := time.Now()
		appId, err := repo.BuildApplication("User", typ, title, remarks, amount, paidAt, []string{"User"})
		if err != nil {
			panic(err)
		}

		newType := ApplicationType{Type: Contest}
		err = repo.PatchApplication(appId, "User", &newType, "", "", nil, nil, []string{})
		asr.NoError(err)

		app, err := repo.GetApplication(appId, true)
		if err != nil {
			panic(err)
		}

		asr.Len(app.ApplicationsDetails, 2)
		asr.Equal(app.LatestApplicationsDetail, app.ApplicationsDetails[1])

		asr.Equal(newType, app.LatestApplicationsDetail.Type)
		asr.Equal(title, app.LatestApplicationsDetail.Title)
		asr.Equal(remarks, app.LatestApplicationsDetail.Remarks)
		asr.Equal(amount, app.LatestApplicationsDetail.Amount)
		asr.Equal(paidAt.Truncate(time.Hour*24), app.LatestApplicationsDetail.PaidAt.PaidAt)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		err = repo.PatchApplication(id, generateRandomUserName(), nil, "", "", nil, nil, []string{})
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})

}

func TestGetApplicationList(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		deleteAllRecord(db)

		user1 := "User1"
		user2 := "User2"
		user3 := "User3"

		app1SubTime := time.Date(2020, 4, 10, 12, 0, 0, 0, time.Local)
		app1Id := buildApplicationWithSubmitTime(user1, app1SubTime, ApplicationType{Type: Club}, "CCCCC", "Remarks", 10000, time.Now())

		app2SubTime := time.Date(2020, 4, 20, 12, 0, 0, 0, time.Local)
		app2Id := buildApplicationWithSubmitTime(user2, app2SubTime, ApplicationType{Type: Contest}, "AAAAA", "Remarks", 10000, time.Now())

		app3SubTime := time.Date(2020, 4, 30, 12, 0, 0, 0, time.Local)
		app3Id := buildApplicationWithSubmitTime(user2, app3SubTime, ApplicationType{Type: Event}, "BBBBB", "Remarks", 10000, time.Now())
		app3, err := repo.GetApplication(app3Id, true)
		if err != nil {
			panic(err)
		}

		app4SubTime := time.Date(2019, 4, 10, 12, 0, 0, 0, time.Local)
		app4Id := buildApplicationWithSubmitTime(user1, app4SubTime, ApplicationType{Type: Club}, "DDDDD", "Remarks", 10000, time.Now())

		// TODO Use a appropriate function defined in model/states_log.go after implementing such a function.
		db.Model(&app3.LatestStatesLog).Updates(StatesLog{
			ToState: StateType{FullyRepaid},
		})

		t.Parallel()

		t.Run("allNil", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetApplicationList("", nil, nil, "", nil, nil, nil)
			asr.NoError(err)

			asr.Len(apps, 4)
			asr.Equal([]uuid.UUID{app3Id, app2Id, app1Id, app4Id}, mapToApplicationID(apps))

			for _, app := range apps {
				asr.NotZero(app.LatestApplicationsDetail)
				asr.NotZero(app.LatestStatesLog)
				asr.NotZero(app.LatestState)
			}

			asr.False(apps[0].CreatedAt.Before(apps[1].CreatedAt))
			asr.False(apps[1].CreatedAt.Before(apps[2].CreatedAt))

			asr.Equal(user1, apps[2].CreateUserTrapID.TrapId)
		})

		t.Run("filterByFinancialYear", func(t *testing.T) {
			asr := assert.New(t)
			financialYear := 2019

			apps, err := repo.GetApplicationList("", nil, &financialYear, "", nil, nil, nil)
			asr.NoError(err)

			asr.Len(apps, 1)
			asr.Equal(app4Id, apps[0].ID)
		})

		t.Run("filterByApplicant", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetApplicationList("", nil, nil, user2, nil, nil, nil)
			asr.NoError(err)

			asr.Len(apps, 2)
		})

		t.Run("filterByCurrentState", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetApplicationList("", &StateType{FullyRepaid}, nil, "", nil, nil, nil)
			asr.NoError(err)

			asr.Len(apps, 1)
			asr.Equal(app3Id, apps[0].ID)
		})

		t.Run("filterByApplicationType", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetApplicationList("", nil, nil, "", &ApplicationType{Type: Contest}, nil, nil)
			asr.NoError(err)

			asr.Len(apps, 1)
			asr.Equal(apps[0].ID, app2Id)
		})

		t.Run("emptyResult", func(t *testing.T) {
			asr := assert.New(t)

			apps, err := repo.GetApplicationList("", nil, nil, user3, nil, nil, nil)
			asr.NoError(err)

			asr.Empty(apps)
		})

		t.Run("filterBySubmitted", func(t *testing.T) {
			beforeApp2 := app2SubTime.Add(-24 * time.Hour)
			beforeApp3 := app3SubTime.Add(-24 * time.Hour)

			t.Parallel()

			t.Run("Since", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := repo.GetApplicationList("", nil, nil, "", nil, &beforeApp3, nil)
				asr.NoError(err)

				asr.Len(apps, 1)
				asr.Equal(app3Id, apps[0].ID)
			})

			t.Run("until", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := repo.GetApplicationList("", nil, nil, "", nil, nil, &beforeApp3)
				asr.NoError(err)

				asr.Len(apps, 3)
				asr.Equal([]uuid.UUID{app2Id, app1Id, app4Id}, mapToApplicationID(apps))
			})

			t.Run("both", func(t *testing.T) {
				asr := assert.New(t)

				apps, err := repo.GetApplicationList("", nil, nil, "", nil, &beforeApp2, &beforeApp3)
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

					apps, err := repo.GetApplicationList(test.SortBy, nil, nil, "", nil, nil, nil)
					asr.NoError(err)

					asr.Len(apps, 4)
					asr.Equal(test.Should, mapToApplicationID(apps))
				})
			}
		})
	})
}

func mapToApplicationID(apps []Application) []uuid.UUID {
	appIds := make([]uuid.UUID, len(apps))
	for i := range apps {
		appIds[i] = apps[i].ID
	}

	return appIds
}
