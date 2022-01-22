package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/transactiondetail"
)

func (repo *EntRepository) GetTransactions(ctx context.Context, query TransactionQuery) ([]*TransactionResponse, error) {
	// Querying
	var transactionsq *ent.TransactionQuery
	var err error
	if query.Sort == nil || *query.Sort == "" || *query.Sort == "created_at" {
		transactionsq = repo.client.Transaction.
			Query().
			WithTag().
			WithDetail().
			WithGroupBudget(func(q *ent.GroupBudgetQuery) {
				q.WithGroup()
			}).
			Order(ent.Desc(transaction.FieldCreatedAt))
	} else if *query.Sort == "-created_at" {
		transactionsq = repo.client.Transaction.
			Query().
			WithTag().
			WithDetail().
			WithGroupBudget(func(q *ent.GroupBudgetQuery) {
				q.WithGroup()
			}).
			Order(ent.Asc(transaction.FieldCreatedAt))
	}

	if query.Target != nil && *query.Target != "" {
		transactionsq = transactionsq.Where(transaction.HasDetailWith(
			transactiondetail.TargetEQ(*query.Target),
		))
	}

	if query.Year != nil && *query.Year != 0 {
		transactionsq = transactionsq.
			Where(transaction.CreatedAtGTE(time.Date(*query.Year, 4, 1, 0, 0, 0, 0, time.Local))).
			Where(transaction.CreatedAtLT(time.Date(*query.Year+1, 4, 1, 0, 0, 0, 0, time.Local)))
	}

	txs, err := transactionsq.All(ctx)
	if err != nil {
		return nil, err
	}

	// Converting
	var res []*TransactionResponse
	for _, tx := range txs {
		res = append(res, ConvertEntTransactionToModelTransactionResponse(tx))
	}

	return res, nil
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

func ConvertEntTransactionToModelTransactionResponse(transaction *ent.Transaction) *TransactionResponse {
	var tags []*Tag
	for _, tag := range transaction.Edges.Tag {
		tags = append(tags, ConvertEntTagToModelTag(tag))
	}
	return &TransactionResponse{
		ID:        transaction.ID,
		Amount:    transaction.Edges.Detail.Amount,
		Target:    transaction.Edges.Detail.Target,
		Tags:      tags,
		Group:     ConvertEntGroupToModelGroup(transaction.Edges.GroupBudget.Edges.Group),
		CreatedAt: transaction.CreatedAt,
	}
}
