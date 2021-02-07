package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

type Request struct {
	ID                  uuid.UUID       `gorm:"type:char(36);primary_key" json:"id"`
	LatestRequestStatus RequestStatus   `gorm:"foreignkey:StatesLogsID" json:"-"`
	LatestStatus        Status          `gorm:"-" json:"current_state"`
	RequestStatusID     int             `gorm:"type:int(11);not null" json:"-"`
	CreatedBy           TrapUser        `gorm:"embedded;embedded_prefix:created_by;" json:"applicant"`
	CreatedAt           time.Time       `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	RequestStatus       []RequestStatus `json:"request_status"`
	Files               []File          `json:"files"`
	Comments            []Comment       `json:"comments"`
	RequestTargets      []RequestTarget `json:"request_targets"`
}

func (app *Request) GiveIsUserAdmin(admins []string) {
	if app == nil {
		return
	}

	app.CreatedBy.GiveIsUserAdmin(admins)
	app.LatestRequestStatus.GiveIsUserAdmin(admins)

	for i := range app.RequestStatus {
		app.RequestStatus[i].GiveIsUserAdmin(admins)
	}

	for i := range app.Comments {
		app.Comments[i].GiveIsUserAdmin(admins)
	}

	for i := range app.RequestTargets {
		app.RequestTargets[i].GiveIsUserAdmin(admins)
	}
}

// RequestRepository 依頼のリポジトリ
type RequestRepository interface {
	GetRequest(id uuid.UUID, preload bool) (Request, error)
	GetRequestList(
		sort string,
		currentState *StateType,
		financialYear *int,
		applicant string,
		submittedSince *time.Time,
		submittedUntil *time.Time,
	) ([]Request, error)
	BuildRequest(
		createdBy string,
		title string,
		remarks string,
		amount int,
		paidAt time.Time,
		repayUsers []string,
	) (uuid.UUID, error)
	PatchRequest(
		appID uuid.UUID,
		updateUserTrapID string,
		title string,
		remarks string,
		amount *int,
		paidAt *time.Time,
		targets []string,
	) error
	UpdateRequestStatus(
		applicationID uuid.UUID,
		updateUserTrapID string,
		reason string,
		status Status,
	) (RequestStatus, error)
	UpdateRequestTarget(
		applicationID uuid.UUID,
		target string,
		repaidAt time.Time,
	) (RequestTarget, bool, error)
}

type requestRepository struct{}

func NewRequestRepository() RequestRepository {
	return &requestRepository{}
}

func (*requestRepository) GetRequest(id uuid.UUID, preload bool) (Request, error) {
	var req Request
	query := db
	if preload {
		query = query.Set("gorm:auto_preload", true)
	} else {
		query = query.Preload("LatestRequestStatus")
	}

	err := query.First(&req, Application{ID: id}).Error
	if err != nil {
		return Request{}, err
	}

	return req, nil
}

func (*requestRepository) GetRequestList(sort string, currentState *StateType, financialYear *int, applicant string, typ *ApplicationType, submittedSince *time.Time, submittedUntil *time.Time) ([]Application, error) {
	query := db.Preload("LatestStatesLog").Preload("LatestApplicationsDetail")

	if currentState != nil {
		query = query.Joins("JOIN states_logs ON states_logs.id = applications.states_logs_id").Where("states_logs.to_state = ?", currentState.Type)
	}

	if financialYear != nil {
		financialYear := time.Date(*financialYear, 4, 1, 0, 0, 0, 0, time.Local)
		financialYearEnd := financialYear.AddDate(1, 0, 0)
		query = query.Where("created_at >= ?", financialYear).Where("created_at < ?", financialYearEnd)
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

func (repo *requestRepository) BuildRequest(createdBy string, title string, remarks string, amount int, paidAt time.Time, targets []string) (uuid.UUID, error) {
	var id uuid.UUID

	err := db.Transaction(func(tx *gorm.DB) error {
		_id, err := repo.createRequest(tx, createdBy)
		if err != nil {
			return err
		}
		id = _id

		state, err := repo.createRequestStatus(tx, id, createdBy)
		if err != nil {
			return err
		}

		for _, userID := range targets {
			if err = repo.createRequestTarget(tx, id, userID); err != nil {
				return err
			}
		}

		return tx.Model(Request{}).Where(&Request{ID: id}).Updates(Request{
			RequestStatusID: state.ID,
		}).Error
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (repo *requestRepository) PatchRequest(appID uuid.UUID, updateUserTrapID string, title string, remarks string, amount *int, paidAt *time.Time, repayUsers []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		app, err := repo.GetRequest(appID, false)
		if err != nil {
			return err
		}

		if len(repayUsers) > 0 {
			if err = repo.deleteRepayUserByApplicationID(tx, appID); err != nil {
				return err
			}

			for _, userID := range repayUsers {
				if err = repo.createRequestTarget(tx, appID, userID); err != nil {
					return err
				}
			}
		}

		return tx.Model(&Application{ID: appID}).Updates(Application{}).Error
	})
}

func (*requestRepository) createRequest(db *gorm.DB, createUserTrapID string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	app := Request{
		ID:        id,
		CreatedBy: TrapUser{TrapID: createUserTrapID},
	}

	err = db.Create(&app).Error
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
