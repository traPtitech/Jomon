package model

import (
	"context"

	"github.com/traPtitech/Jomon/ent"
)

func (repo *EntRepository) CreateTransactionDetail(ctx context.Context, amount int, target string) (*ent.TransactionDetail, error) {
	return repo.client.TransactionDetail.Create().
		SetAmount(amount).
		SetTarget(target).
		Save(ctx)
}
