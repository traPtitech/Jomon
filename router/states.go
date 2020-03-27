package router

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/Jomon/model"

	"github.com/labstack/echo/v4"
)

type StatesLog struct {
	ID               int       `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID    uuid.UUID `gorm:"type:char(36);not null" json:"-"`
	UpdateUserTrapID model.User      `gorm:"embedded;embedded_prefix:update_user_" json:"update_user"`
	ToState          model.StateType `gorm:"embedded" json:"to_state"`
	Reason           string    `gorm:"type:text;not null" json:"reason"`
	CreatedAt        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type ChangeState struct {
	ToState 		model.StateType `gorm:"embedded" json:"to_state"`
	Reason 			string `gorm:"type:text;not null" json:"reason"`
}

type ReturnState struct {
	User			model.User `gorm:"embedded;embedded_prefix:user" json:"user"`
	UpdatedAt 		time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CurrentState 	model.StateType `gorm:"embedded" json:"current_state"`
	PastState 		model.StateType `gorm:"embedded" json:"past_state"`
}

type DebugState struct {
	CurrentState model.StateType `gorm:"embedded" json:"current_state"`
	ToState model.StateType `gorm:"embedded" json:"to_state"`
}

func (s *Service) PutStates(c echo.Context) error {
	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	} // applicationIdがとれた
	application, err := s.Applications.GetApplication(applicationId, true)
	if gorm.IsRecordNotFoundError(err) {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	currentState := application.LatestStatesLog.ToState
	updateState := new(ChangeState)
	_err := c.Bind(updateState)
	if _err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	toState := updateState.ToState
	debugState := DebugState{
		CurrentState: currentState,
		ToState: toState,
	}
	return c.JSON(http.StatusOK, debugState)
}

func (s *Service) PutRepaidStates(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutRepaidStates")
}
