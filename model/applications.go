package model

import (
	"time"

	"github.com/labstack/echo/v4"
)

type SmallApplication struct {
	ApplicationID string    `json:"application_id"`
	CreatedAt     time.Time `json:"created_at"`
	Applicant     Applicant `json:"applicant"`
	CurrentState  string    `json:"current_state"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	Remarks       string    `json:"remarks"`
	PaidAt        string    `json:"paid_at"`
	Ammount       int       `json:"ammount"`
}

type Applicant struct {
	TrapID      string `json:"trap_id"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
}

type Application struct {
	ApplicationID         string                `json:"application_id"`
	CreatedAt             time.Time             `json:"created_at"`
	Applicant             Applicant             `json:"applicant"`
	CurrentState          string                `json:"current_state"`
	Type                  string                `json:"type"`
	Title                 string                `json:"title"`
	Remarks               string                `json:"remarks"`
	PaidAt                string                `json:"paid_at"`
	Ammount               int                   `json:"ammount"`
	RepaidToID            []string              `json:"repaid_to_id"`
	Images                []string              `json:"images"`
	Comments              Comments              `json:"comments"`
	StateLogs             StateLogs             `json:"state_logs"`
	ApplicationDetailLogs ApplicationDetailLogs `json:"application_detail_logs"`
	RepaymentLogs         RepaymentLogs         `json:"repayment_logs"`
}

type Comments []struct {
	CommentID int       `json:"comment_id"`
	User      User      `json:"user"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	TrapID      string `json:"trap_id"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
}

type StateLogs []struct {
	UpdateUser UpdateUser `json:"update_user"`
	ToState    string     `json:"to_state"`
	Reason     string     `json:"reason"`
	CreatedAt  time.Time  `json:"created_at"`
}

type UpdateUser struct {
	TrapID      string `json:"trap_id"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
}

type ApplicationDetailLogs []struct {
	UpdateUser UpdateUser `json:"update_user"`
	Type       string     `json:"type"`
	Title      string     `json:"title"`
	Remarks    string     `json:"remarks"`
	Ammount    int        `json:"ammount"`
	PaidAt     string     `json:"paid_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type RepaymentLogs []struct {
	RepaidByUser RepaidByUser `json:"repaid_by_user"`
	RepaidToUser RepaidToUser `json:"repaid_to_user"`
	RepaidAt     time.Time    `json:"repaid_at"`
}

type RepaidByUser struct {
	TrapID      string `json:"trap_id"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
}

type RepaidToUser struct {
	TrapID      string `json:"trap_id"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
}

func GetApplications(c echo.Context) ([]SmallApplication, error) {
	allapplications := []SmallApplication{}
	return allapplications, nil
}

func PostApplications(c echo.Context) (Application, error) {
	application := Application{}
	return application, nil
}

func PatchApplications(c echo.Context) (Application, error) {
	application := Application{}
	return application, nil
}
