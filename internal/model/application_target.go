package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/ent/application"
	"github.com/traPtitech/Jomon/internal/ent/applicationtarget"
	"github.com/traPtitech/Jomon/internal/nulltime"
	"github.com/traPtitech/Jomon/internal/service"
)

func (repo *EntRepository) GetApplicationTargets(
	ctx context.Context, applicationID uuid.UUID,
) ([]*service.ApplicationTargetDetail, error) {
	errorConverter := &entErrorConverter{
		msgBadInput: "failed to get application targets due to invalid input",
		msgNotFound: "application targets not found",
	}
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
		return nil, errorConverter.convert(err)
	}
	targets := lo.Map(ts, func(t *ent.ApplicationTarget, _ int) *service.ApplicationTargetDetail {
		return ConvertEntApplicationTargetToModelApplicationTargetDetail(t)
	})
	return targets, nil
}

func (repo *EntRepository) createApplicationTargets(
	ctx context.Context, tx *ent.Tx, applicationID uuid.UUID, targets []*service.ApplicationTarget,
) ([]*service.ApplicationTargetDetail, error) {
	bulk := lo.Map(targets, func(t *service.ApplicationTarget, _ int) *ent.ApplicationTargetCreate {
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
	ts := lo.Map(created, func(t *ent.ApplicationTarget, _ int) *service.ApplicationTargetDetail {
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

func ConvertEntApplicationTargetToModelApplicationTargetDetail(
	t *ent.ApplicationTarget,
) *service.ApplicationTargetDetail {
	return &service.ApplicationTargetDetail{
		ID:        t.ID,
		Target:    t.Edges.User.ID,
		Amount:    t.Amount,
		PaidAt:    nulltime.FromTime(t.PaidAt),
		CreatedAt: t.CreatedAt,
	}
}
