package router

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/traPtitech/Jomon/model"

	"github.com/labstack/echo/v4"
)

type PutState struct {
	ToState model.StateType `json:"to_state"`
	Reason  string          `json:"reason"`
}

type SuccessState struct {
	User         model.User      `json:"user"`
	UpdatedAt    time.Time       `json:"updated_at"`
	CurrentState model.StateType `json:"current_state"`
	PastState    model.StateType `json:"past_state"`
}

type ErrorState struct {
	CurrentState model.StateType `json:"current_state"`
	ToState      model.StateType `json:"to_state"`
}

type SuccessRepaid struct {
	RepaidByUser model.User      `json:"repaid_by_user_trap_id"`
	RepaidToUser model.User      `json:"repaid_to_user_trap_id"`
	RepaidAt     *time.Time      `json:"repaid_at"`
	ToState      model.StateType `json:"to_state"`
}

func (s *Service) PutStates(c echo.Context) error {
	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	application, err := s.Applications.GetApplication(applicationId, false)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c.NoContent(http.StatusNotFound)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	var sta PutState
	if err := c.Bind(&sta); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	errsta := &ErrorState{
		CurrentState: application.LatestState,
		ToState:      sta.ToState,
	}

	if sta.Reason == "" {
		if IsAbleNoReasonChangeState(sta.ToState, application.LatestState) {
			return c.JSON(http.StatusBadRequest, errsta)
		}
	}

	if user == application.CreateUserTrapID {
		if IsAbleCreatorChangeState(sta.ToState, application.LatestState) {
			return c.JSON(http.StatusBadRequest, errsta)
		}
	}

	admin, err := s.Administrators.IsAdmin(user.TrapId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	if admin {
		if IsAbleAdminChangeState(sta.ToState, application.LatestState) {
			return c.JSON(http.StatusBadRequest, errsta)
		}
	}

	state, err := s.Applications.UpdateStatesLog(applicationId, user.TrapId, sta.Reason, sta.ToState)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	sucsta := &SuccessState{
		User:         user,
		UpdatedAt:    state.CreatedAt,
		CurrentState: state.ToState,
		PastState:    application.LatestState,
	}

	return c.JSON(http.StatusOK, sucsta)
}

func IsAbleNoReasonChangeState(toState model.StateType, currentState model.StateType) bool {
	if (toState == model.StateType{Type: model.FixRequired}) && (currentState == model.StateType{Type: model.Submitted}) {
		return true
	}
	if (toState == model.StateType{Type: model.Rejected}) && (currentState == model.StateType{Type: model.Submitted}) {
		return true
	}
	if (toState == model.StateType{Type: model.Submitted}) && (currentState == model.StateType{Type: model.Accepted}) {
		return true
	}
	return false
}
func IsAbleCreatorChangeState(toState model.StateType, currentState model.StateType) bool {
	if (toState == model.StateType{Type: model.Submitted}) && (currentState == model.StateType{Type: model.FixRequired}) {
		return false
	}
	return true
}

func IsAbleAdminChangeState(toState model.StateType, currentState model.StateType) bool {
	if (toState == model.StateType{Type: model.Rejected}) && (currentState == model.StateType{Type: model.Submitted}) {
		return false
	}
	if (toState == model.StateType{Type: model.FixRequired}) && (currentState == model.StateType{Type: model.Submitted}) {
		return false
	}
	if (toState == model.StateType{Type: model.Submitted}) && (currentState == model.StateType{Type: model.FixRequired}) {
		return false
	}
	if (toState == model.StateType{Type: model.Accepted}) && (currentState == model.StateType{Type: model.Submitted}) {
		return false
	}
	if (toState == model.StateType{Type: model.Submitted}) && (currentState == model.StateType{Type: model.Accepted}) {
		return false
	}
	return true
}

func (s *Service) PutRepaidStates(c echo.Context) error {
	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	repaidToId := c.Param("repaidToId")
	if repaidToId == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	application, err := s.Applications.GetApplication(applicationId, false)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c.NoContent(http.StatusNotFound)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	admin, err := s.Administrators.IsAdmin(user.TrapId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	if (!admin) || (application.LatestState != model.StateType{Type: model.Accepted}) {
		return c.NoContent(http.StatusForbidden)
	}

	updateRepayUser, allUsersRepaidCheck, err := s.Applications.UpdateRepayUser(applicationId, repaidToId, user.TrapId)
	switch {
	case err == model.ErrAlreadyRepaid:
		return c.NoContent(http.StatusBadRequest)
	case err != nil:
		return c.NoContent(http.StatusInternalServerError)
	}

	var sucrep *SuccessRepaid
	if allUsersRepaidCheck {
		sucrep = &SuccessRepaid{
			RepaidByUser: model.User{
				TrapId: updateRepayUser.RepaidByUserTrapID.TrapId,
			},
			RepaidToUser: model.User{
				TrapId: updateRepayUser.RepaidToUserTrapID.TrapId,
			},
			RepaidAt: updateRepayUser.RepaidAt,
			ToState:  model.StateType{Type: model.FullyRepaid},
		}
	} else {
		sucrep = &SuccessRepaid{
			RepaidByUser: model.User{
				TrapId: updateRepayUser.RepaidByUserTrapID.TrapId,
			},
			RepaidToUser: model.User{
				TrapId: updateRepayUser.RepaidToUserTrapID.TrapId,
			},
			RepaidAt: updateRepayUser.RepaidAt,
			ToState:  model.StateType{Type: model.Submitted},
		}
	}

	return c.JSON(http.StatusOK, sucrep)
}
