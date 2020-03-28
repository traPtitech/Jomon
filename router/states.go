package router

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/Jomon/model"

	"github.com/labstack/echo/v4"
)

type PutState struct {
	ToState 		model.StateType `gorm:"embedded" json:"to_state"`
	Reason 			string `gorm:"type:text;not null" json:"reason"`
}

type SuccessState struct {
	User model.User `gorm:"embedded" json:"user"`
	UpdatedAt time.Time `gorm:"embedded" json:"updated_at"`
	CurrentState model.StateType `gorm:"embedded" json:"current_state"`
	PastState model.StateType `gorm:"embedded" json:"past_state"`
}

type ErrorState struct {
	CurrentState 	model.StateType `gorm:"embedded" json:"current_state"`
	ToState 		model.StateType `gorm:"embedded" json:"to_state"`
}

func (s *Service) PutStates(c echo.Context) error {
	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	application, err := s.Applications.GetApplication(applicationId, false)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.NoContent(http.StatusNotFound)
		} else {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	var sta PutState
	if err := c.Bind(sta); err != nil{
		return c.NoContent(http.StatusBadRequest)
	}

	errsta := &ErrorState{
		CurrentState: application.LatestState,
		ToState: sta.ToState,
	}

	if sta.Reason == "" {
		if IsNoReasonState(sta.ToState, application.LatestState) {
			return c.JSON(http.StatusBadRequest, errsta)
		}
	}

	if user == application.CreateUserTrapID {
		if IsCreatorState(sta.ToState, application.LatestState) {
			return c.JSON(http.StatusBadRequest, errsta)
		}
	}

	admin, err := s.Administrators.IsAdmin(user.TrapId)
	if err != nil{
		return c.NoContent(http.StatusBadRequest)
	}
	if admin {
		if IsAdminState(sta.ToState, application.LatestState) {
			return c.JSON(http.StatusBadRequest, errsta)
		}
	}
	
	state, err := s.States.CreateStatesLog(applicationId, user.TrapId, sta.Reason, sta.ToState)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	sucsta := &SuccessState{
		User: user,
		UpdatedAt: state.CreatedAt,
		CurrentState: state.ToState,
		PastState: application.LatestState,
	}

	return c.JSON(http.StatusOK, sucsta)
}

func IsNoReasonState(toState model.StateType, currentState model.StateType) bool {
	if (toState == model.StateType{Type: 2}) && (currentState == model.StateType{Type: 1}) {
		return true
	}
	if (toState == model.StateType{Type: 5}) && (currentState == model.StateType{Type: 1}) {
		return true
	}
	if (toState == model.StateType{Type: 1}) && (currentState == model.StateType{Type: 3}) {
		return true
	}
	return false
}

func IsCreatorState(toState model.StateType, currentState model.StateType) bool {
	if (toState == model.StateType{Type: 1}) && (currentState == model.StateType{Type: 2}) {
		return false
	}
	return true
}

func IsAdminState(toState model.StateType, currentState model.StateType) bool {
	if (toState == model.StateType{Type: 5}) && (currentState == model.StateType{Type: 1}) {
		return false
	}
	if (toState == model.StateType{Type: 2}) && (currentState == model.StateType{Type: 1}) {
		return false
	}
	if (toState == model.StateType{Type: 1}) && (currentState == model.StateType{Type: 2}) {
		return false
	}
	if (toState == model.StateType{Type: 3}) && (currentState == model.StateType{Type: 1}) {
		return false
	}
	if (toState == model.StateType{Type: 1}) && (currentState == model.StateType{Type: 3}) {
		return false
	}
	return true
}

func (s *Service) PutRepaidStates(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutRepaidStates")
}
