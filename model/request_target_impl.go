package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requesttarget"
)

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

func createRequestTargets(client *ent.Client, ctx context.Context, requestID uuid.UUID, targets []*Target) ([]*TargetDetail, error) {
	var entTargets []*ent.RequestTarget
	bulk := make([]*ent.RequestTargetCreate, len(targets))
	for i, target := range targets {
		bulk[i] = client.RequestTarget.Create().SetTarget(target.Target).SetAmount(target.Amount).SetRequestID(requestID)
	}
	entTargets, err := client.RequestTarget.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, err
	}
	var targetDetails []*TargetDetail
	for _, entTarget := range entTargets {
		targetDetails = append(targetDetails, convertEntRequestTargetToModelRequestTarget(entTarget))
	}
	return targetDetails, nil
}

func updateRequestTargets(client *ent.Client, ctx context.Context, requestID uuid.UUID, targets []Target) ([]*TargetDetail, error) {
	var entTargets []*ent.RequestTarget
	bulk := make([]*ent.RequestTargetCreate, len(targets))
	for i, target := range targets {
		bulk[i] = client.RequestTarget.Create().SetTarget(target.Target).SetAmount(target.Amount).SetRequestID(requestID)
	}
	entTargets, err := client.RequestTarget.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, err
	}
	var targetDetails []*TargetDetail
	for _, entTarget := range entTargets {
		targetDetails = append(targetDetails, convertEntRequestTargetToModelRequestTarget(entTarget))
	}
	return targetDetails, nil
}

func deleteRequestTargets(client *ent.Client, ctx context.Context, requestID uuid.UUID) error {
	_, err := client.RequestTarget.
		Delete().
		Where(requesttarget.HasRequestWith(request.IDEQ(requestID))).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func convertEntRequestTargetToModelRequestTarget(requestTarget *ent.RequestTarget) *TargetDetail {
	if requestTarget == nil {
		return nil
	}
	return &TargetDetail{
		ID:        requestTarget.ID,
		Target:    requestTarget.Target,
		Amount:    requestTarget.Amount,
		PaidAt:    requestTarget.PaidAt,
		CreatedAt: requestTarget.CreatedAt,
	}
}
