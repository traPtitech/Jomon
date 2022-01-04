package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
)

func (repo *EntRepository) GetTransactions(ctx context.Context, query TransactionQuery) ([]*TransactionResponse, error) {
	// TODO: impl
	return nil, nil
}

func (repo *EntRepository) GetTransaction(ctx context.Context, transactionID uuid.UUID) (*TransactionResponse, error) {
	// TODO: impl
	return nil, nil
}

func (repo *EntRepository) CreateTransaction(ctx context.Context, Amount int, Target string, tags []*uuid.UUID, group *uuid.UUID) (*TransactionResponse, error) {
	// TODO: impl
	return nil, nil
}

func (repo *EntRepository) UpdateTransaction(ctx context.Context, transactionID uuid.UUID, Amount int, Target string, tags []*uuid.UUID, group *uuid.UUID) (*TransactionResponse, error) {
	// TODO: impl
	return nil, nil
}

func ConvertEntTransactionToModelTransaction(transaction *ent.Transaction) *Transaction {
	return &Transaction{
		ID:        transaction.ID,
		CreatedAt: transaction.CreatedAt,
	}
}
