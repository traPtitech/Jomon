package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/transactiondetail"
)

func (repo *EntRepository) CreateTransactionDetail(ctx context.Context, amount int, target string) (*TransactionDetail, error) {
	enttd, err := repo.client.TransactionDetail.Create().
		SetAmount(amount).
		SetTarget(target).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return convertEntTransactionDetailToModelTransactionDetail(enttd), nil
}

func (repo *EntRepository) UpdateTransactionDetail(ctx context.Context, transactionID uuid.UUID, amount int, target string) (*TransactionDetail, error) {
	_, err := repo.client.TransactionDetail.Update().
		Where(transactiondetail.HasTransactionWith(
			transaction.IDEQ(transactionID),
		)).
		SetAmount(amount).
		SetTarget(target).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Get Transaction Detail
	enttd, err := repo.client.TransactionDetail.
		Query().
		Where(transactiondetail.HasTransactionWith(
			transaction.IDEQ(transactionID),
		)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return convertEntTransactionDetailToModelTransactionDetail(enttd), nil
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
