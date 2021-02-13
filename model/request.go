package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

// Request struct of Request
type Request struct {
	ID                  uuid.UUID       `gorm:"type:char(36);primary_key" json:"id"`
	LatestRequestStatus RequestStatus   `gorm:"foreignkey:RequestStatusID" json:"-"`
	LatestStatus        Status          `gorm:"-" json:"current_state"`
	RequestStatusID     int             `gorm:"type:int(11);not null" json:"-"`
	Amount              int             `gorm:"type:int(11);not null" json:"amount"`
	Title               string          `gorm:"type:text;not null" json:"title"`
	Content             string          `gorm:"type:text;not null" json:"content"`
	CreatedBy           TrapUser        `gorm:"embedded;embedded_prefix:created_by_;" json:"applicant"`
	CreatedAt           time.Time       `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	RequestStatus       []RequestStatus `json:"request_status"`
	Files               []File          `json:"files"`
	Comments            []Comment       `json:"comments"`
	RequestTargets      []RequestTarget `json:"request_targets"`
}

// GiveIsUserAdmin check whether request is admin or not
func (req *Request) GiveIsUserAdmin(admins []string) {
	if req == nil {
		return
	}

	req.CreatedBy.GiveIsUserAdmin(admins)
	req.LatestRequestStatus.GiveIsUserAdmin(admins)

	for i := range req.RequestStatus {
		req.RequestStatus[i].GiveIsUserAdmin(admins)
	}

	for i := range req.Comments {
		req.Comments[i].GiveIsUserAdmin(admins)
	}
}

// RequestRepository 依頼のリポジトリ
type RequestRepository interface {
	GetRequest(id uuid.UUID, preload bool) (Request, error)
	GetRequestList(sort string,
		currentStatus *Status,
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
		createdBy string,
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
		requestID uuid.UUID,
		target string,
		createdBy string,
		paidAt time.Time,
	) (RequestTarget, bool, error)
}

type requestRepository struct{}

// NewRequestRepository Make RequestRepository
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

	err := query.First(&req, Request{ID: id}).Error
	if err != nil {
		return Request{}, err
	}

	return req, nil
}

func (*requestRepository) GetRequestList(sort string, currentStatus *Status, financialYear *int, applicant string, submittedSince *time.Time, submittedUntil *time.Time) ([]Request, error) {
	query := db.Preload("LatestRequestStatus")

	if currentStatus != nil {
		query = query.Joins("JOIN request_status ON request_status.request_id = request.id").Where("request_status.status = ?", currentStatus)
	}

	if financialYear != nil {
		financialYear := time.Date(*financialYear, 4, 1, 0, 0, 0, 0, time.Local)
		financialYearEnd := financialYear.AddDate(1, 0, 0)
		query = query.Where("created_at >= ?", financialYear).Where("created_at < ?", financialYearEnd)
	}

	if applicant != "" {
		query = query.Where("created_by = ?", applicant)
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
		query = query.Order("title")
	case "-title":
		query = query.Order("title desc")
	}

	//noinspection GoPreferNilSlice
	apps := []Request{}
	err := query.Find(&apps).Error
	if err != nil {
		return nil, err
	}

	for i := range apps {
		apps[i].LatestStatus = apps[i].LatestRequestStatus.Status
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

func (repo *requestRepository) PatchRequest(reqID uuid.UUID, createdBy string, title string, content string, amount *int, paidAt *time.Time, requestTargets []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		req, err := repo.GetRequest(reqID, false)
		if err != nil {
			return err
		}

		req.ID, err = uuid.NewV4() // zero value of int is 0
		if err != nil {
			return err
		}

		req.CreatedBy.TrapID = createdBy

		if title != "" {
			req.Title = title
		}

		if content != "" {
			req.Content = content
		}

		if amount != nil {
			req.Amount = *amount
		}

		req.UpdatedAt = time.Time{} // zero value

		if len(requestTargets) > 0 {
			if err = repo.deleteRequestTargetByRequestID(tx, reqID); err != nil {
				return err
			}

			for _, userID := range requestTargets {
				if err = repo.createRequestTarget(tx, reqID, userID); err != nil {
					return err
				}
			}
		}

		return tx.Model(&Request{ID: reqID}).Updates(Request{}).Error
	})
}

func (*requestRepository) createRequest(db *gorm.DB, createdBy string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	app := Request{
		ID:        id,
		CreatedBy: TrapUser{TrapID: createdBy},
	}

	err = db.Create(&app).Error
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
