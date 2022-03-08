package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
)

func (repo *EntRepository) CreateTransactionDetail(ctx context.Context, amount int, target string) (*TransactionDetail, error) {
	enttd, err := repo.client.TransactionDetail.Create().
		SetAmount(amount).
		SetTarget(target).
		Save(ctx)
	return convertEntTransactionDetailToModelTransactionDetail(enttd), err
}

func (repo *EntRepository) UpdateTransactionDetail(ctx context.Context, transactionID uuid.UUID, amount int, target string) (*TransactionDetail, error) {
	enttd, err := repo.client.TransactionDetail.UpdateOneID(transactionID).
		SetAmount(amount).
		SetTarget(target).
		Save(ctx)
	return convertEntTransactionDetailToModelTransactionDetail(enttd), err
}

func convertEntTransactionDetailToModelTransactionDetail(enttd *ent.TransactionDetail) *TransactionDetail {
	if enttd == nil {
		return nil
	}
	return &TransactionDetail{
		ID:        enttd.ID,
		Amount:    enttd.Amount,
		Target:    enttd.Target,
		CreatedAt: enttd.CreatedAt,
		UpdatedAt: enttd.UpdatedAt,
	}
}
