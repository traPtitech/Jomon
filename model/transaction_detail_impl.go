package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/transactiondetail"
)

func (repo *EntRepository) createTransactionDetail(
	ctx context.Context, tx *ent.Tx, title string, amount int, target string,
) (*TransactionDetail, error) {
	enttd, err := tx.Client().TransactionDetail.Create().
		SetTitle(title).
		SetAmount(amount).
		SetTarget(target).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return convertEntTransactionDetailToModelTransactionDetail(enttd), nil
}

func (repo *EntRepository) updateTransactionDetail(
	ctx context.Context, tx *ent.Tx, transactionID uuid.UUID,
	title string, amount int, target string,
) (*TransactionDetail, error) {
	_, err := tx.Client().TransactionDetail.Update().
		Where(transactiondetail.HasTransactionWith(
			transaction.IDEQ(transactionID),
		)).
		ClearTransaction().
		Save(ctx)
	if err != nil {
		return nil, err
	}
	enttd, err := tx.Client().TransactionDetail.Create().
		SetTitle(title).
		SetAmount(amount).
		SetTarget(target).
		SetTransactionID(transactionID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return convertEntTransactionDetailToModelTransactionDetail(enttd), nil
}

func convertEntTransactionDetailToModelTransactionDetail(
	enttd *ent.TransactionDetail,
) *TransactionDetail {
	if enttd == nil {
		return nil
	}
	return &TransactionDetail{
		ID:        enttd.ID,
		Title:     enttd.Title,
		Amount:    enttd.Amount,
		Target:    enttd.Target,
		CreatedAt: enttd.CreatedAt,
		UpdatedAt: enttd.UpdatedAt,
	}
}
