package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requesttarget"
)

func (repo *EntRepository) GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*RequestTargetDetail, error) {
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
	var targets []*RequestTargetDetail
	for _, t := range ts {
		targets = append(targets, ConvertEntRequestTargetToModelRequestTargetDetail(t))
	}
	return targets, err
}

func (repo *EntRepository) createRequestTargets(ctx context.Context, tx *ent.Tx, requestID uuid.UUID, targets []*RequestTarget) ([]*RequestTargetDetail, error) {
	var bulk []*ent.RequestTargetCreate
	for _, t := range targets {
		bulk = append(bulk,
			tx.Client().RequestTarget.Create().
				SetAmount(t.Amount).
				SetRequestID(requestID).
				SetUserID(t.Target),
		)
	}
	cs, err := tx.Client().RequestTarget.CreateBulk(bulk...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	ids := []uuid.UUID{}
	for _, c := range cs {
		ids = append(ids, c.ID)
	}
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
	var ts []*RequestTargetDetail
	for _, t := range created {
		ts = append(ts, ConvertEntRequestTargetToModelRequestTargetDetail(t))
	}
	return ts, nil
}

func (repo *EntRepository) deleteRequestTargets(ctx context.Context, tx *ent.Tx, requestID uuid.UUID) error {
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
	fmt.Printf("hoge %#v\n", t)
	return &RequestTargetDetail{
		ID:        t.ID,
		Target:    t.Edges.User.ID,
		Amount:    t.Amount,
		PaidAt:    t.PaidAt,
		CreatedAt: t.CreatedAt,
	}
}
