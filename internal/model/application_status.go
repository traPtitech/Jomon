package model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/ent/applicationstatus"
	"github.com/traPtitech/Jomon/internal/service"
)

func (repo *EntRepository) CreateStatus(
	ctx context.Context, applicationID uuid.UUID, userID uuid.UUID, status service.Status,
) (*service.ApplicationStatus, error) {
	errorConverter := &entErrorConverter{
		msgBadInput: "failed to create application status due to invalid input",
		msgNotFound: "application status not found",
	}
	c, err := repo.client.ApplicationStatus.
		Create().
		SetStatus(applicationstatus.Status(status.String())).
		SetCreatedAt(time.Now()).
		SetApplicationID(applicationID).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, errorConverter.convert(err)
	}
	created, err := repo.client.ApplicationStatus.
		Query().
		Where(applicationstatus.ID(c.ID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, errorConverter.convert(err)
	}
	return convertEntApplicationStatusToModelApplicationStatus(created), nil
}

func convertEntApplicationStatusToModelApplicationStatus(
	applicationStatus *ent.ApplicationStatus,
) *service.ApplicationStatus {
	if applicationStatus == nil {
		return nil
	}
	return &service.ApplicationStatus{
		ID:        applicationStatus.ID,
		CreatedBy: applicationStatus.Edges.User.ID,
		Status:    convertEntApplicationStatusToModelStatus(&applicationStatus.Status),
		CreatedAt: applicationStatus.CreatedAt,
	}
}

func convertEntApplicationStatusToModelStatus(entStatus *applicationstatus.Status) service.Status {
	var status service.Status
	switch entStatus.String() {
	case "submitted":
		status = service.Submitted
	case "fix_required":
		status = service.FixRequired
	case "accepted":
		status = service.Accepted
	case "completed":
		status = service.Completed
	case "rejected":
		status = service.Rejected
	default:
		panic(fmt.Sprintf("unknown application status: %s", entStatus.String()))
	}
	return status
}
