package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/requeststatus"
)

func (repo *EntRepository) CreateStatus(
	ctx context.Context, requestID uuid.UUID, userID uuid.UUID, status Status,
) (*RequestStatus, error) {
	c, err := repo.client.RequestStatus.
		Create().
		SetStatus(requeststatus.Status(status.String())).
		SetCreatedAt(time.Now()).
		SetRequestID(requestID).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	created, err := repo.client.RequestStatus.
		Query().
		Where(requeststatus.ID(c.ID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return convertEntRequestStatusToModelRequestStatus(created), nil
}

func convertEntRequestStatusToModelRequestStatus(requestStatus *ent.RequestStatus) *RequestStatus {
	if requestStatus == nil {
		return nil
	}
	return &RequestStatus{
		ID:        requestStatus.ID,
		CreatedBy: requestStatus.Edges.User.ID,
		Status:    convertEntRequestStatusToModelStatus(&requestStatus.Status),
		CreatedAt: requestStatus.CreatedAt,
	}
}

func convertEntRequestStatusToModelStatus(entStatus *requeststatus.Status) Status {
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
