package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/applicationstatus"
)

func (repo *EntRepository) CreateStatus(
	ctx context.Context, applicationID uuid.UUID, userID uuid.UUID, status Status,
) (*ApplicationStatus, error) {
	c, err := repo.client.ApplicationStatus.
		Create().
		SetStatus(applicationstatus.Status(status.String())).
		SetCreatedAt(time.Now()).
		SetApplicationID(applicationID).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	created, err := repo.client.ApplicationStatus.
		Query().
		Where(applicationstatus.ID(c.ID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return convertEntApplicationStatusToModelApplicationStatus(created), nil
}

func convertEntApplicationStatusToModelApplicationStatus(
	applicationStatus *ent.ApplicationStatus,
) *ApplicationStatus {
	if applicationStatus == nil {
		return nil
	}
	return &ApplicationStatus{
		ID:        applicationStatus.ID,
		CreatedBy: applicationStatus.Edges.User.ID,
		Status:    convertEntApplicationStatusToModelStatus(&applicationStatus.Status),
		CreatedAt: applicationStatus.CreatedAt,
	}
}

func convertEntApplicationStatusToModelStatus(entStatus *applicationstatus.Status) Status {
	var status Status
	switch entStatus.String() {
	case "submitted":
		status = Submitted
	case "fix_required":
		status = FixRequired
	case "accepted":
		status = Accepted
	case "completed":
		status = Completed
	case "rejected":
		status = Rejected
	}
	return status
}
