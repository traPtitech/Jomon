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
	GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*RequestTarget, error)
}

type RequestTarget struct {
	ID        uuid.UUID
	Target    string
	PaidAt    *time.Time
	CreatedAt time.Time
}

func (repo *EntRepository) GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*RequestTarget, error) {
	targets, err := repo.client.RequestTarget.
		Query().
		Where(requesttarget.IDEQ(requestID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var reqTargets []*RequestTarget
	for _, target := range targets {
		reqTargets = append(reqTargets, convertEntRequestTargetToModelRequestTarget(target))
	}
	return reqTargets, nil
}

func convertEntRequestTargetToModelRequestTarget(requestTarget *ent.RequestTarget) *RequestTarget {
	if requestTarget == nil {
		return nil
	}
	return &RequestTarget{
		ID:        requestTarget.ID,
		Target:    requestTarget.Target,
		PaidAt:    requestTarget.PaidAt,
		CreatedAt: requestTarget.CreatedAt,
	}
}
