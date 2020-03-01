package model

import (
	"github.com/jinzhu/gorm"
	"time"

	"github.com/gofrs/uuid"
)

type Application struct {
	ID                       uuid.UUID            `gorm:"type:char(36);primary_key" json:"application_id"`
	LatestApplicationsDetail ApplicationsDetail   `gorm:"foreignkey:ApplicationsDetailsID" json:"current_detail"`
	ApplicationsDetailsID    int                  `gorm:"type:int(11);not null" json:"-"`
	LatestStatesLog          StatesLog            `gorm:"foreignkey:StatesLogsID" json:"-"`
	LatestStatus             StateType            `gorm:"-" json:"current_state"`
	StatesLogsID             int                  `gorm:"type:int(11);not null" json:"-"`
	CreateUserTrapID         User                 `gorm:"embedded;embedded_prefix:create_user_;" json:"applicant"`
	CreatedAt                time.Time            `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
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
	GetApplication(id uuid.UUID, giveAdmin bool, preload bool) (Application, error)
	GetApplicationList(
		sort *string,
		currentState *StateType,
		financialYear *int,
		applicant *string,
		typ *ApplicationType,
		submittedSince *time.Time,
		submittedUntil *time.Time,
		giveAdmin bool,
	) ([]Application, error)
	BuildApplication(
		createUserTrapID string,
		typ ApplicationType,
		title string,
		remarks string,
		amount int,
		paidAt time.Time,
	) (uuid.UUID, error)
	PatchApplication(
		appId uuid.UUID,
		updateUserTrapId string,
		typ *ApplicationType,
		title *string,
		remarks *string,
		amount *int,
		paidAt *time.Time,
	) error
}

type applicationRepository struct{}

func NewApplicationRepository() ApplicationRepository {
	return &applicationRepository{}
}

func (_ *applicationRepository) GetApplication(id uuid.UUID, giveAdmin bool, preload bool) (Application, error) {
	var app Application
	query := db
	if preload {
		query = query.Set("gorm:auto_preload", true)
	}

	err := query.First(&app, Application{ID: id}).Error
	if err != nil {
		return Application{}, err
	}

	app.LatestStatus = app.LatestStatesLog.ToState
	if giveAdmin {
		admins, err := GetAdministratorList()
		if err != nil {
			return Application{}, err
		}
		app.GiveIsUserAdmin(admins)
	}

	return app, nil
}

func (_ *applicationRepository) GetApplicationList(sort *string, currentState *StateType, financialYear *int, applicant *string, typ *ApplicationType, submittedSince *time.Time, submittedUntil *time.Time, giveAdmin bool) ([]Application, error) {
	query := db

	if currentState != nil {
		query = query.Joins("JOIN states_logs ON states_logs.id = applications.states_logs_id").Where("states_logs.to_state = ?", currentState.Type)
	}

	if applicant != nil {
		query = query.Where("create_user_trap_id = ?", *applicant)
	}

	if typ != nil {
		query = query.Joins("JOIN applications_details ON applications_details.id = applications.applications_details_id").Where("applications_details.type = ?", typ.Type)
	}

	if submittedSince != nil {
		query = query.Where("created_at > ?", *submittedSince)
	}

	if submittedUntil != nil {
		query = query.Where("created_at < ?", *submittedUntil)
	}

	if sort != nil {
		switch *sort {
		case "created_at":
			query = query.Order("created_at desc")
		case "-created_at":
			query = query.Order("created_at")
		case "title":
			query = query.Joins("JOIN applications_details ON applications_details.id = applications.applications_details_id").Order("applications_details.title")
		case "-title":
			query = query.Joins("JOIN applications_details ON applications_details.id = applications.applications_details_id").Order("applications_details.title desc")
		}
	} else {
		query = query.Order("created_at desc")
	}

	//noinspection GoPreferNilSlice
	apps := []Application{}
	err := query.Find(&apps).Error
	if err != nil {
		return nil, err
	}

	for i := range apps {
		apps[i].LatestStatus = apps[i].LatestStatesLog.ToState
	}

	if giveAdmin {
		admins, err := GetAdministratorList()
		if err != nil {
			return nil, err
		}

		for i := range apps {
			apps[i].GiveIsUserAdmin(admins)
		}
	}

	return apps, nil
}

func (repo *applicationRepository) BuildApplication(createUserTrapID string, typ ApplicationType, title string, remarks string, amount int, paidAt time.Time) (uuid.UUID, error) {
	var id uuid.UUID

	err := db.Transaction(func(tx *gorm.DB) error {
		_id, err := repo.createApplication(tx, createUserTrapID)
		if err != nil {
			return err
		}
		id = _id

		detail, err := createApplicationsDetail(tx, id, createUserTrapID, typ, title, remarks, amount, paidAt)
		if err != nil {
			return err
		}

		state, err := createStatesLog(tx, id, createUserTrapID)
		if err != nil {
			return err
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

func (repo *applicationRepository) PatchApplication(appId uuid.UUID, updateUserTrapId string, typ *ApplicationType, title *string, remarks *string, amount *int, paidAt *time.Time) error {
	return db.Transaction(func(tx *gorm.DB) error {
		app, err := repo.GetApplication(appId, false, false)
		if err != nil {
			return err
		}

		detail, err := putApplicationsDetail(tx, app.ApplicationsDetailsID, updateUserTrapId, typ, title, remarks, amount, paidAt)
		if err != nil {
			return err
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
