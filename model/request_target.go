//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/requesttarget"
)

type RequestTargetRepository interface {
	GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*TargetDetail, error)
}

type Target struct {
	Target string
	Amount int
}
type TargetDetail struct {
	ID        uuid.UUID
	Target    string
	Amount    int
	PaidAt    *time.Time
	CreatedAt time.Time
}

func (repo *EntRepository) GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*TargetDetail, error) {
	targets, err := repo.client.RequestTarget.
		Query().
		Where(requesttarget.IDEQ(requestID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var reqTargets []*TargetDetail
	for _, target := range targets {
		reqTargets = append(reqTargets, convertEntRequestTargetToModelRequestTarget(target))
	}
	return reqTargets, nil
}

func convertEntRequestTargetToModelRequestTarget(requestTarget *ent.RequestTarget) *TargetDetail {
	if requestTarget == nil {
		return nil
	}
	return &TargetDetail{
		ID:        requestTarget.ID,
		Target:    requestTarget.Target,
		PaidAt:    requestTarget.PaidAt,
		CreatedAt: requestTarget.CreatedAt,
	}
}
