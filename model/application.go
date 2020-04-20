package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

type Application struct {
	ID                       uuid.UUID            `gorm:"type:char(36);primary_key" json:"application_id"`
	LatestApplicationsDetail ApplicationsDetail   `gorm:"foreignkey:ApplicationsDetailsID" json:"current_detail"`
	ApplicationsDetailsID    int                  `gorm:"type:int(11);not null" json:"-"`
	LatestStatesLog          StatesLog            `gorm:"foreignkey:StatesLogsID" json:"-"`
	LatestState              StateType            `gorm:"-" json:"current_state"`
	StatesLogsID             int                  `gorm:"type:int(11);not null" json:"-"`
	CreateUserTrapID         User                 `gorm:"embedded;embedded_prefix:create_user_;" json:"applicant"`
	CreatedAt                time.Time            `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ApplicationsDetails      []ApplicationsDetail `json:"application_detail_logs,omitempty"`
	StatesLogs               []StatesLog          `json:"state_logs,omitempty"`
	ApplicationsImages       []ApplicationsImage  `json:"images,omitempty"`
	Comments                 []Comment            `json:"comments,omitempty"`
	RepayUsers               []RepayUser          `json:"repayment_logs,omitempty" `
}

func (app *Application) GiveIsUserAdmin(admins []string) {
	if app == nil {
		return
	}

	app.CreateUserTrapID.GiveIsUserAdmin(admins)
	app.LatestApplicationsDetail.GiveIsUserAdmin(admins)
	app.LatestStatesLog.GiveIsUserAdmin(admins)

	for i := range app.ApplicationsDetails {
		app.ApplicationsDetails[i].GiveIsUserAdmin(admins)
	}

	for i := range app.StatesLogs {
		app.StatesLogs[i].GiveIsUserAdmin(admins)
	}

	for i := range app.Comments {
		app.Comments[i].GiveIsUserAdmin(admins)
	}

	for i := range app.RepayUsers {
		app.RepayUsers[i].GiveIsUserAdmin(admins)
	}
}

type ApplicationRepository interface {
	GetApplication(id uuid.UUID, preload bool) (Application, error)
	GetApplicationList(
		sort string,
		currentState *StateType,
		financialYear *int,
		applicant string,
		typ *ApplicationType,
		submittedSince *time.Time,
		submittedUntil *time.Time,
	) ([]Application, error)
	BuildApplication(
		createUserTrapID string,
		typ ApplicationType,
		title string,
		remarks string,
		amount int,
		paidAt time.Time,
		repayUsers []string,
	) (uuid.UUID, error)
	PatchApplication(
		appId uuid.UUID,
		updateUserTrapId string,
		typ *ApplicationType,
		title string,
		remarks string,
		amount *int,
		paidAt *time.Time,
		repayUsers []string,
	) error
	UpdateStatesLog(
		applicationId uuid.UUID,
		updateUserTrapId string,
		reason string,
		toState StateType,
	) (StatesLog, error)
	UpdateRepayUser(
		applicationId uuid.UUID,
		repaidToUserTrapID string,
		repaidByUserTrapID string,
	) (RepayUser, bool, error)
}

type applicationRepository struct{}

func NewApplicationRepository() ApplicationRepository {
	return &applicationRepository{}
}

func (_ *applicationRepository) GetApplication(id uuid.UUID, preload bool) (Application, error) {
	var app Application
	query := db
	if preload {
		query = query.Set("gorm:auto_preload", true)
	} else {
		query = query.Preload("LatestStatesLog")
	}

	err := query.First(&app, Application{ID: id}).Error
	if err != nil {
		return Application{}, err
	}

	app.LatestState = app.LatestStatesLog.ToState

	return app, nil
}

//noinspection GoUnusedParameter `financialYear`
func (_ *applicationRepository) GetApplicationList(sort string, currentState *StateType, financialYear *int, applicant string, typ *ApplicationType, submittedSince *time.Time, submittedUntil *time.Time) ([]Application, error) {
	query := db.Preload("LatestStatesLog").Preload("LatestApplicationsDetail")

	if currentState != nil {
		query = query.Joins("JOIN states_logs ON states_logs.id = applications.states_logs_id").Where("states_logs.to_state = ?", currentState.Type)
	}

	if applicant != "" {
		query = query.Where("create_user_trap_id = ?", applicant)
	}

	if typ != nil {
		query = query.Joins("JOIN applications_details ON applications_details.id = applications.applications_details_id").Where("applications_details.type = ?", typ.Type)
	}

	if submittedSince != nil {
		query = query.Where("created_at >= ?", *submittedSince)
	}

	if submittedUntil != nil {
		query = query.Where("created_at < ?", *submittedUntil)
	}

	switch sort {
	case "", "created_at":
		query = query.Order("created_at desc")
	case "-created_at":
		query = query.Order("created_at")
	case "title":
		query = query.Joins("JOIN applications_details ON applications_details.id = applications.applications_details_id").Order("applications_details.title")
	case "-title":
		query = query.Joins("JOIN applications_details ON applications_details.id = applications.applications_details_id").Order("applications_details.title desc")
	}

	//noinspection GoPreferNilSlice
	apps := []Application{}
	err := query.Find(&apps).Error
	if err != nil {
		return nil, err
	}

	for i := range apps {
		apps[i].LatestState = apps[i].LatestStatesLog.ToState
	}

	return apps, nil
}

func (repo *applicationRepository) BuildApplication(createUserTrapID string, typ ApplicationType, title string, remarks string, amount int, paidAt time.Time, repayUsers []string) (uuid.UUID, error) {
	var id uuid.UUID

	err := db.Transaction(func(tx *gorm.DB) error {
		_id, err := repo.createApplication(tx, createUserTrapID)
		if err != nil {
			return err
		}
		id = _id

		detail, err := repo.createApplicationsDetail(tx, id, createUserTrapID, typ, title, remarks, amount, paidAt)
		if err != nil {
			return err
		}

		state, err := repo.createStatesLog(tx, id, createUserTrapID)
		if err != nil {
			return err
		}

		for _, userId := range repayUsers {
			if err = repo.createRepayUser(tx, id, userId); err != nil {
				return err
			}
		}

		return tx.Model(Application{}).Where(&Application{ID: id}).Updates(Application{
			ApplicationsDetailsID: detail.ID,
			StatesLogsID:          state.ID,
		}).Error
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (repo *applicationRepository) PatchApplication(appId uuid.UUID, updateUserTrapId string, typ *ApplicationType, title string, remarks string, amount *int, paidAt *time.Time, repayUsers []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		app, err := repo.GetApplication(appId, false)
		if err != nil {
			return err
		}

		detail, err := repo.putApplicationsDetail(tx, app.ApplicationsDetailsID, updateUserTrapId, typ, title, remarks, amount, paidAt)
		if err != nil {
			return err
		}

		if len(repayUsers) > 0 {
			if err = repo.deleteRepayUserByApplicationID(tx, appId); err != nil {
				return err
			}

			for _, userId := range repayUsers {
				if err = repo.createRepayUser(tx, appId, userId); err != nil {
					return err
				}
			}
		}

		return tx.Model(&Application{ID: appId}).Updates(Application{
			ApplicationsDetailsID: detail.ID,
		}).Error
	})
}

func (_ *applicationRepository) createApplication(db_ *gorm.DB, createUserTrapID string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	app := Application{
		ID:               id,
		CreateUserTrapID: User{TrapId: createUserTrapID},
	}

	err = db_.Create(&app).Error
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
