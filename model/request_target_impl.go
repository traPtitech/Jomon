package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requesttarget"
	"github.com/traPtitech/Jomon/service"
)

func (repo *EntRepository) GetRequestTargets(
	ctx context.Context, requestID uuid.UUID,
) ([]*RequestTargetDetail, error) {
	// Querying
	ts, err := repo.client.RequestTarget.
		Query().
		Where(
			requesttarget.HasRequestWith(
				request.IDEQ(requestID),
			),
		).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	targets := lo.Map(ts, func(t *ent.RequestTarget, _ int) *RequestTargetDetail {
		return ConvertEntRequestTargetToModelRequestTargetDetail(t)
	})
	return targets, err
}

func (repo *EntRepository) createRequestTargets(
	ctx context.Context, tx *ent.Tx, requestID uuid.UUID, targets []*RequestTarget,
) ([]*RequestTargetDetail, error) {
	bulk := lo.Map(targets, func(t *RequestTarget, _ int) *ent.RequestTargetCreate {
		return tx.Client().RequestTarget.
			Create().
			SetAmount(t.Amount).
			SetRequestID(requestID).
			SetUserID(t.Target)
	})
	cs, err := tx.Client().RequestTarget.CreateBulk(bulk...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	ids := lo.Map(cs, func(c *ent.RequestTarget, _ int) uuid.UUID {
		return c.ID
	})
	created, err := tx.Client().RequestTarget.
		Query().
		Where(
			requesttarget.IDIn(ids...),
		).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	// []*ent.RequestTarget to []*RequestTargetDetail
	ts := lo.Map(created, func(t *ent.RequestTarget, _ int) *RequestTargetDetail {
		return ConvertEntRequestTargetToModelRequestTargetDetail(t)
	})
	return ts, nil
}

func (repo *EntRepository) deleteRequestTargets(
	ctx context.Context, tx *ent.Tx, requestID uuid.UUID,
) error {
	_, err := tx.Client().RequestTarget.
		Delete().
		Where(
			requesttarget.HasRequestWith(
				request.IDEQ(requestID),
			),
		).
		Exec(ctx)
	return err
}

func ConvertEntRequestTargetToModelRequestTargetDetail(t *ent.RequestTarget) *RequestTargetDetail {
	return &RequestTargetDetail{
		ID:        t.ID,
		Target:    t.Edges.User.ID,
		Amount:    t.Amount,
		PaidAt:    service.TimeToNullTime(t.PaidAt).Time,
		CreatedAt: t.CreatedAt,
	}
}
