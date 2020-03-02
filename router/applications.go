package router

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/traPtitech/Jomon/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type GetApplicationsListQuery struct {
	Sort           string `query:"sort"`
	CurrentState   string `query:"current_state"`
	FinancialYear  *int   `query:"financial_year"`
	Applicant      string `query:"applicant"`
	Type           string `query:"type"`
	SubmittedSince string `query:"submitted_since"`
	SubmittedUntil string `query:"submitted_until"`
}

type PostApplicationRequest struct {
	Type       *model.ApplicationType `json:"type"`
	Title      string                 `json:"title"`
	Remarks    string                 `json:"remarks"`
	PaidAt     *time.Time             `json:"paid_at"`
	Amount     *int                   `json:"amount"`
	RepaidToId []string               `json:"repaid_to_id"`
}

func (s *Service) GetApplicationList(c echo.Context) error {
	var query GetApplicationsListQuery
	err := c.Bind(&query)
	if err != nil {
		// TODO
		return c.NoContent(http.StatusBadRequest)
	}

	var currentState *model.StateType
	if query.CurrentState != "" {
		_currentState, err := model.GetStateType(query.CurrentState)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		currentState = &_currentState
	}

	var typ *model.ApplicationType
	if query.Type != "" {
		_typ, err := model.GetApplicationType(query.Type)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		typ = &_typ
	}

	var submittedSince *time.Time
	if query.SubmittedSince != "" {
		_submittedSince, err := StrToDate(query.SubmittedSince)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		submittedSince = &_submittedSince
	}

	var submittedUntil *time.Time
	if query.SubmittedUntil != "" {
		_submittedUntil, err := StrToDate(query.SubmittedUntil)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		_submittedUntil = _submittedUntil.AddDate(0, 0, 1)
		submittedUntil = &_submittedUntil
	}

	applications, err := s.Applications.GetApplicationList(query.Sort, currentState, query.FinancialYear, query.Applicant, typ, submittedSince, submittedUntil)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	for i := range applications {
		applications[i].GiveIsUserAdmin(admins)
	}

	return c.JSON(http.StatusOK, applications)
}

func (s *Service) GetApplication(c echo.Context) error {
	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}
	application, err := s.Applications.GetApplication(applicationId, true)
	if gorm.IsRecordNotFoundError(err) {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	application.GiveIsUserAdmin(admins)

	return c.JSON(http.StatusOK, application)
}

func (s *Service) PostApplication(c echo.Context) error {
	details := c.FormValue("details")
	var req PostApplicationRequest
	if err := json.Unmarshal([]byte(details), &req); err != nil {
		// TODO more information
		return c.NoContent(http.StatusBadRequest)
	}

	if req.Type == nil || req.Title == "" || req.Remarks == "" || req.Amount == nil || req.PaidAt == nil || len(req.RepaidToId) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	userId := "UserId"

	id, err := s.Applications.BuildApplication(userId, *req.Type, req.Title, req.Remarks, *req.Amount, *req.PaidAt, req.RepaidToId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	app, err := s.Applications.GetApplication(id, true)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	app.GiveIsUserAdmin(admins)

	// TODO Handle image files

	return c.JSON(http.StatusCreated, app)
}

type PatchErrorMessage struct {
	IsUserAccepted bool             `json:"is_user_accepted"`
	CurrentState   *model.StateType `json:"current_state,omitempty"`
}

func (s *Service) PatchApplication(c echo.Context) error {
	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	userId := "UserId"
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	details := c.FormValue("details")
	var req PostApplicationRequest
	if err := json.Unmarshal([]byte(details), &req); err != nil {
		// TODO
		return c.NoContent(http.StatusBadRequest)
	}

	if req.Type == nil && req.Title == "" && req.Remarks == "" && req.Amount == nil && req.PaidAt == nil && (len(req.RepaidToId) == 0) {
		return c.NoContent(http.StatusBadRequest)
	}

	app, err := s.Applications.GetApplication(applicationId, true)
	isRequestUserAdmin, err := s.Administrators.IsAdmin(userId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	} else if !isRequestUserAdmin {
		if app.CreateUserTrapID.TrapId != userId {
			return c.JSON(http.StatusForbidden, &PatchErrorMessage{
				IsUserAccepted: false,
			})
		} else if app.LatestState.Type != model.Submitted && app.LatestState.Type != model.FixRequired {
			return c.JSON(http.StatusForbidden, &PatchErrorMessage{
				IsUserAccepted: true,
				CurrentState:   &app.LatestState,
			})
		}
	}

	err = s.Applications.PatchApplication(applicationId, userId, req.Type, req.Title, req.Remarks, req.Amount, req.PaidAt, req.RepaidToId)
	if gorm.IsRecordNotFoundError(err) {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	app, err = s.Applications.GetApplication(applicationId, true)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	admins, err := s.Administrators.GetAdministratorList()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	app.GiveIsUserAdmin(admins)

	// TODO Handle image files

	return c.JSON(http.StatusOK, &app)
}
