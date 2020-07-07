package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/traPtitech/Jomon/model"

	"github.com/labstack/echo/v4"
)

var acceptedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/bmp":  true,
}

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
	form, err := c.MultipartForm()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	details := form.Value["details"][0]
	var req PostApplicationRequest
	if err := json.Unmarshal([]byte(details), &req); err != nil {
		// TODO more information
		return c.NoContent(http.StatusBadRequest)
	}

	if req.Type == nil || req.Title == "" || req.Remarks == "" || req.Amount == nil || req.PaidAt == nil || len(req.RepaidToId) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	if *req.Amount <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	user, ok := c.Get(contextUserKey).(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	id, err := s.Applications.BuildApplication(user.TrapId, *req.Type, req.Title, req.Remarks, *req.Amount, *req.PaidAt, req.RepaidToId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	images := form.File["images"]
	for _, file := range images {
		mimeType := file.Header.Get(echo.HeaderContentType)
		if !acceptedMimeTypes[mimeType] {
			return c.NoContent(http.StatusUnsupportedMediaType)
		}

		src, err := file.Open()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		defer src.Close()

		_, err = s.Images.CreateApplicationsImage(id, src, mimeType)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
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

	return c.JSON(http.StatusCreated, app)
}

type PatchErrorMessage struct {
	IsUserAccepted bool             `json:"is_user_accepted"`
	CurrentState   *model.StateType `json:"current_state,omitempty"`
}

func (s *Service) PatchApplication(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	details := form.Value["details"][0]
	var req PostApplicationRequest
	if err := json.Unmarshal([]byte(details), &req); err != nil {
		// TODO
		return c.NoContent(http.StatusBadRequest)
	}

	if req.Type == nil && req.Title == "" && req.Remarks == "" && req.Amount == nil && req.PaidAt == nil && (len(req.RepaidToId) == 0) && (len(form.File["images"]) == 0) {
		return c.NoContent(http.StatusBadRequest)
	}

	if req.Amount != nil && *req.Amount <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	app, err := s.Applications.GetApplication(applicationId, true)
	if gorm.IsRecordNotFoundError(err) {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user, ok := c.Get(contextUserKey).(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	if !user.IsAdmin {
		if app.CreateUserTrapID.TrapId != user.TrapId {
			return c.JSON(http.StatusForbidden, &PatchErrorMessage{
				IsUserAccepted: false,
			})
		} else if app.LatestState.Type != model.Submitted && app.LatestState.Type != model.FixRequired {
			return c.JSON(http.StatusForbidden, &PatchErrorMessage{
				IsUserAccepted: true,
				CurrentState:   &app.LatestState,
			})
		}
	} else if app.LatestState.Type != model.Submitted && app.LatestState.Type != model.FixRequired {
		return c.JSON(http.StatusForbidden, &PatchErrorMessage{
			IsUserAccepted: true,
			CurrentState:   &app.LatestState,
		})
	}

	err = s.Applications.PatchApplication(applicationId, user.TrapId, req.Type, req.Title, req.Remarks, req.Amount, req.PaidAt, req.RepaidToId)
	if gorm.IsRecordNotFoundError(err) {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	images := form.File["images"]
	for _, file := range images {
		mimeType := file.Header.Get(echo.HeaderContentType)
		if !acceptedMimeTypes[mimeType] {
			return c.NoContent(http.StatusUnsupportedMediaType)
		}

		src, err := file.Open()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		defer src.Close()

		_, err = s.Images.CreateApplicationsImage(applicationId, src, mimeType)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
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

	return c.JSON(http.StatusOK, &app)
}
