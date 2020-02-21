package model

import (
	"encoding/json"
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

func (app Application) MarshalJSON() ([]byte, error) {
	app.LatestStatus = app.LatestStatesLog.ToState
	return json.Marshal(app)
}

func (app *Application) GiveIsUserAdmin(admins []string) {
	if app == nil {
		return
	}

	app.CreateUserTrapID.GiveIsUserAdmin(admins)
	app.LatestApplicationsDetail.GiveIsUserAdmin(admins)

	for _, det := range app.ApplicationsDetails {
		det.GiveIsUserAdmin(admins)
	}

	for _, st := range app.StatesLogs {
		st.GiveIsUserAdmin(admins)
	}

	for _, com := range app.Comments {
		com.GiveIsUserAdmin(admins)
	}

	for _, ru := range app.RepayUsers {
		ru.GiveIsUserAdmin(admins)
	}
}

func GetApplication(id uuid.UUID, giveAdmin bool) (Application, error) {
	var app Application
	err := db.Set("gorm:auto_preload", true).First(&app, Application{ID: id}).Error
	if err != nil {
		return Application{}, err
	}

	if giveAdmin {
		admins, err := GetAdministratorList()
		if err != nil {
			return Application{}, err
		}
		app.GiveIsUserAdmin(admins)
	}

	return Application{}, nil
}

func GetApplicationList(sort *string, currentState *StateType, financialYear *int, applicant *string, typ *ApplicationType, submittedSince *time.Time, submittedUntil *time.Time, giveAdmin bool) ([]Application, error) {
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
	}

	//noinspection GoPreferNilSlice
	apps := []Application{}
	err := query.Find(&apps).Error
	if err != nil {
		return nil, err
	}

	if giveAdmin {
		admins, err := GetAdministratorList()
		if err != nil {
			return nil, err
		}

		for _, app := range apps {
			app.GiveIsUserAdmin(admins)
		}
	}

	return apps, nil
}

func CreateApplication(createUserTrapID string, giveAdmin bool) (Application, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Application{}, err
	}

	app := Application{
		ID:               id,
		CreateUserTrapID: User{TrapId: createUserTrapID},
	}
	db.Create(&app)

	if giveAdmin {
		admins, err := GetAdministratorList()
		if err != nil {
			return Application{}, err
		}
		app.GiveIsUserAdmin(admins)
	}

	return Application{}, nil
}

func UpdateApplicationDetailId(id uuid.UUID, newApplicationDetailId int) (Application, error) {
	app := Application{
		ID: id,
	}

	err := db.Model(&app).Updates(Application{
		ApplicationsDetailsID: newApplicationDetailId,
	}).Error
	if err != nil {
		return Application{}, err
	}

	return Application{}, nil
}

func UpdateStatesLogsId(id uuid.UUID, newStatesLogsId int) (Application, error) {
	app := Application{
		ID: id,
	}

	err := db.Model(&app).Updates(Application{
		StatesLogsID: newStatesLogsId,
	}).Error
	if err != nil {
		return Application{}, err
	}

	return Application{}, nil
}
