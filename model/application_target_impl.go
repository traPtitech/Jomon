package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/application"
	"github.com/traPtitech/Jomon/ent/applicationtarget"
)

func (repo *EntRepository) GetApplicationTargets(
	ctx context.Context, applicationID uuid.UUID,
) ([]*ApplicationTargetDetail, error) {
	// Querying
	ts, err := repo.client.ApplicationTarget.
		Query().
		Where(
			applicationtarget.HasApplicationWith(
				application.IDEQ(applicationID),
			),
		).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	targets := lo.Map(ts, func(t *ent.ApplicationTarget, _ int) *ApplicationTargetDetail {
		return ConvertEntApplicationTargetToModelApplicationTargetDetail(t)
	})
	return targets, err
}

func (repo *EntRepository) createApplicationTargets(
	ctx context.Context, tx *ent.Tx, applicationID uuid.UUID, targets []*ApplicationTarget,
) ([]*ApplicationTargetDetail, error) {
	bulk := lo.Map(targets, func(t *ApplicationTarget, _ int) *ent.ApplicationTargetCreate {
		return tx.Client().ApplicationTarget.
			Create().
			SetAmount(t.Amount).
			SetApplicationID(applicationID).
			SetUserID(t.Target)
	})
	cs, err := tx.Client().ApplicationTarget.CreateBulk(bulk...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	ids := lo.Map(cs, func(c *ent.ApplicationTarget, _ int) uuid.UUID {
		return c.ID
	})
	created, err := tx.Client().ApplicationTarget.
		Query().
		Where(
			applicationtarget.IDIn(ids...),
		).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	// []*ent.ApplicationTarget to []*ApplicationTargetDetail
	ts := lo.Map(created, func(t *ent.ApplicationTarget, _ int) *ApplicationTargetDetail {
		return ConvertEntApplicationTargetToModelApplicationTargetDetail(t)
	})
	return ts, nil
}

func (repo *EntRepository) deleteApplicationTargets(
	ctx context.Context, tx *ent.Tx, applicationID uuid.UUID,
) error {
	_, err := tx.Client().ApplicationTarget.
		Delete().
		Where(
			applicationtarget.HasApplicationWith(
				application.IDEQ(applicationID),
			),
		).
		Exec(ctx)
	return err
}

func ConvertEntApplicationTargetToModelApplicationTargetDetail(t *ent.ApplicationTarget) *ApplicationTargetDetail {
	return &ApplicationTargetDetail{
		ID:        t.ID,
		Target:    t.Edges.User.ID,
		Amount:    t.Amount,
		PaidAt:    t.PaidAt,
		CreatedAt: t.CreatedAt,
	}
}
