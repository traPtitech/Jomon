package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/group"
	"github.com/traPtitech/Jomon/ent/groupbudget"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/tag"
	"github.com/traPtitech/Jomon/ent/transaction"
	"github.com/traPtitech/Jomon/ent/transactiondetail"
)

func (repo *EntRepository) GetTransactions(ctx context.Context, query TransactionQuery) ([]*TransactionResponse, error) {
	// Querying
	var transactionsq *ent.TransactionQuery
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

	if query.Since != nil {
		transactionsq = transactionsq.
			Where(transaction.CreatedAtGTE(*query.Since))
	}

	if query.Until != nil {
		transactionsq = transactionsq.
			Where(transaction.CreatedAtLT(*query.Until))
	}

	if query.Tag != nil {
		transactionsq = transactionsq.
			Where(transaction.HasTagWith(
				tag.NameEQ(*query.Tag),
			))
	}

	if query.Group != nil {
		transactionsq = transactionsq.
			Where(transaction.HasGroupBudgetWith(
				groupbudget.HasGroupWith(
					group.NameEQ(*query.Group),
				),
			))
	}

	if query.Request != nil {
		transactionsq = transactionsq.
			Where(transaction.HasRequestWith(
				request.IDEQ(*query.Request),
			))
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
	// Querying
	tx, err := repo.client.Transaction.
		Query().
		Where(transaction.ID(transactionID)).
		WithTag().
		WithDetail().
		WithGroupBudget(func(q *ent.GroupBudgetQuery) {
			q.WithGroup()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Converting
	return ConvertEntTransactionToModelTransactionResponse(tx), nil
}

func (repo *EntRepository) CreateTransaction(ctx context.Context, amount int, target string, tags []*uuid.UUID, groupID *uuid.UUID, requestID *uuid.UUID) (*TransactionResponse, error) {
	// Get Tags
	var tagIDs []uuid.UUID
	for _, tag := range tags {
		tagIDs = append(tagIDs, *tag)
	}

	// Create Transaction Detail
	detail, err := repo.CreateTransactionDetail(ctx, amount, target)
	if err != nil {
		return nil, err
	}

	// Create GroupBudget
	var gb *ent.GroupBudget
	if groupID != nil {
		gb, err = repo.client.GroupBudget.
			Create().
			SetGroupID(*groupID).
			SetAmount(amount).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}
	// Create transaction query
	query := repo.client.Transaction.
		Create().
		SetDetailID(detail.ID)

	// Set Query of GroupBudget if exists
	if gb != nil {
		query.
			SetGroupBudgetID(gb.ID)
	}

	// Set Request to the Transaction
	if requestID != nil {
		query.
			SetRequestID(*requestID)
	}

	// Create transaction
	tx, err := query.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Update Tag to set transaction
	if tags != nil {
		_, err = repo.client.Tag.
			Update().
			Where(tag.IDIn(tagIDs...)).
			AddTransactionIDs(tx.ID).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	tx, err = repo.client.Transaction.
		Query().
		Where(transaction.ID(tx.ID)).
		WithTag().
		WithDetail().
		WithGroupBudget(func(q *ent.GroupBudgetQuery) {
			q.WithGroup()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Converting
	return ConvertEntTransactionToModelTransactionResponse(tx), nil
}

func (repo *EntRepository) UpdateTransaction(ctx context.Context, transactionID uuid.UUID, amount int, target string, tags []*uuid.UUID, groupID *uuid.UUID, requestID *uuid.UUID) (*TransactionResponse, error) {
	// Update transaction Detail
	_, err := repo.UpdateTransactionDetail(ctx, transactionID, amount, target)
	if err != nil {
		return nil, err
	}

	// Get Tags
	var tagIDs []uuid.UUID
	for _, tag := range tags {
		tagIDs = append(tagIDs, *tag)
	}

	// Delete Tag Transaction Edge
	_, err = repo.client.Tag.
		Update().
		Where(tag.HasTransactionWith(
			transaction.IDEQ(transactionID),
		)).
		RemoveTransactionIDs(transactionID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if groupID != nil {
		// Delete GroupBudget
		_, err = repo.client.GroupBudget.
			Delete().
			Where(groupbudget.HasTransactionWith(
				transaction.IDEQ(transactionID),
			)).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		// Create GroupBudget
		_, err = repo.client.GroupBudget.
			Create().
			SetGroupID(*groupID).
			SetAmount(amount).
			SetTransactionID(transactionID).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Update Tag to set transaction
	_, err = repo.client.Tag.
		Update().
		Where(tag.IDIn(tagIDs...)).
		AddTransactionIDs(transactionID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Get transaction
	tx, err := repo.client.Transaction.
		Query().
		Where(transaction.ID(transactionID)).
		WithTag().
		WithDetail().
		WithGroupBudget(func(q *ent.GroupBudgetQuery) {
			q.WithGroup()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Converting
	return ConvertEntTransactionToModelTransactionResponse(tx), nil
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
	var group *Group
	if transaction.Edges.GroupBudget == nil {
		group = nil
	} else {
		group = ConvertEntGroupToModelGroup(transaction.Edges.GroupBudget.Edges.Group)
	}
	return &TransactionResponse{
		ID:        transaction.ID,
		Amount:    transaction.Edges.Detail.Amount,
		Target:    transaction.Edges.Detail.Target,
		Tags:      tags,
		Group:     group,
		CreatedAt: transaction.CreatedAt,
	}
}
