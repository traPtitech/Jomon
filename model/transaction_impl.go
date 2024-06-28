package model

import (
	"context"

	"entgo.io/ent/dialect/sql"
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
			WithRequest().
			WithGroupBudget(func(q *ent.GroupBudgetQuery) {
				q.WithGroup()
			}).
			Order(ent.Desc(transaction.FieldCreatedAt))
	} else if *query.Sort == "-created_at" {
		transactionsq = repo.client.Transaction.
			Query().
			WithTag().
			WithDetail().
			WithRequest().
			WithGroupBudget(func(q *ent.GroupBudgetQuery) {
				q.WithGroup()
			}).
			Order(ent.Asc(transaction.FieldCreatedAt))
	} else if *query.Sort == "amount" {
		transactionsq = repo.client.Transaction.
			Query().
			WithTag().
			WithDetail().
			WithRequest().
			WithGroupBudget(func(q *ent.GroupBudgetQuery) {
				q.WithGroup()
			}).
			Order(func(s *sql.Selector) {
				t := sql.Table(transactiondetail.Table)
				s.Join(t).On(s.C(transaction.FieldID), t.C(transactiondetail.TransactionColumn))
				s.OrderBy(sql.Asc(t.C(transactiondetail.FieldAmount)))
			})
	} else if *query.Sort == "-amount" {
		transactionsq = repo.client.Transaction.
			Query().
			WithTag().
			WithDetail().
			WithRequest().
			WithGroupBudget(func(q *ent.GroupBudgetQuery) {
				q.WithGroup()
			}).
			Order(func(s *sql.Selector) {
				t := sql.Table(transactiondetail.Table)
				s.Join(t).On(s.C(transaction.FieldID), t.C(transactiondetail.TransactionColumn))
				s.OrderBy(sql.Desc(t.C(transactiondetail.FieldAmount)))
			})
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

	transactionsq = transactionsq.Limit(query.Limit).Offset(query.Offset)

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
		WithRequest().
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
	tx, err := repo.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// Get Tags
	var tagIDs []uuid.UUID
	for _, t := range tags {
		tagIDs = append(tagIDs, *t)
	}

	// Create Transaction Detail
	detail, err := repo.createTransactionDetail(ctx, tx, amount, target)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	// Create GroupBudget
	var gb *ent.GroupBudget
	if groupID != nil {
		gb, err = tx.Client().GroupBudget.
			Create().
			SetGroupID(*groupID).
			SetAmount(amount).
			Save(ctx)
		if err != nil {
			err = RollbackWithError(tx, err)
			return nil, err
		}
	}
	// Create transaction query
	query := tx.Client().Transaction.
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
	trns, err := query.Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	// Update Tag to set transaction
	if tags != nil {
		_, err = tx.Client().Tag.
			Update().
			Where(tag.IDIn(tagIDs...)).
			AddTransactionIDs(trns.ID).
			Save(ctx)
		if err != nil {
			err = RollbackWithError(tx, err)
			return nil, err
		}
	}

	trns, err = tx.Client().Transaction.
		Query().
		Where(transaction.ID(trns.ID)).
		WithTag().
		WithDetail().
		WithRequest().
		WithGroupBudget(func(q *ent.GroupBudgetQuery) {
			q.WithGroup()
		}).
		Only(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Converting
	return ConvertEntTransactionToModelTransactionResponse(trns), nil
}

func (repo *EntRepository) UpdateTransaction(ctx context.Context, transactionID uuid.UUID, amount int, target string, tags []*uuid.UUID, groupID *uuid.UUID, requestID *uuid.UUID) (*TransactionResponse, error) {
	tx, err := repo.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// Update transaction Detail
	_, err = repo.updateTransactionDetail(ctx, tx, transactionID, amount, target)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	// Get Tags
	var tagIDs []uuid.UUID
	for _, t := range tags {
		tagIDs = append(tagIDs, *t)
	}

	// Delete Tag Transaction Edge
	_, err = tx.Client().Tag.
		Update().
		Where(tag.HasTransactionWith(
			transaction.IDEQ(transactionID),
		)).
		RemoveTransactionIDs(transactionID).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	if groupID != nil {
		// Delete GroupBudget
		_, err = tx.Client().GroupBudget.
			Delete().
			Where(groupbudget.HasTransactionWith(
				transaction.IDEQ(transactionID),
			)).
			Exec(ctx)
		if err != nil {
			err = RollbackWithError(tx, err)
			return nil, err
		}
		// Create GroupBudget
		_, err = tx.Client().GroupBudget.
			Create().
			SetGroupID(*groupID).
			SetAmount(amount).
			AddTransactionIDs(transactionID).
			Save(ctx)
		if err != nil {
			err = RollbackWithError(tx, err)
			return nil, err
		}
	}

	// Update Tag to set transaction
	_, err = tx.Client().Tag.
		Update().
		Where(tag.IDIn(tagIDs...)).
		AddTransactionIDs(transactionID).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	// Update Request to the Transaction
	if requestID != nil {
		_, err = tx.Client().Transaction.
			UpdateOneID(transactionID).
			SetRequestID(*requestID).
			Save(ctx)
		if err != nil {
			err = RollbackWithError(tx, err)
			return nil, err
		}
	} else {
		_, err = tx.Client().Transaction.
			UpdateOneID(transactionID).
			ClearRequest().
			Save(ctx)
		if err != nil {
			err = RollbackWithError(tx, err)
			return nil, err
		}
	}

	// Get transaction
	trns, err := tx.Client().Transaction.
		Query().
		Where(transaction.ID(transactionID)).
		WithTag().
		WithDetail().
		WithRequest().
		WithGroupBudget(func(q *ent.GroupBudgetQuery) {
			q.WithGroup()
		}).
		Only(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Converting
	return ConvertEntTransactionToModelTransactionResponse(trns), nil
}

func ConvertEntTransactionToModelTransaction(transaction *ent.Transaction) *Transaction {
	return &Transaction{
		ID:        transaction.ID,
		CreatedAt: transaction.CreatedAt,
	}
}

func ConvertEntTransactionToModelTransactionResponse(transaction *ent.Transaction) *TransactionResponse {
	var tags []*Tag
	for _, t := range transaction.Edges.Tag {
		tags = append(tags, ConvertEntTagToModelTag(t))
	}
	var g *Group
	if transaction.Edges.GroupBudget != nil {
		g = ConvertEntGroupToModelGroup(transaction.Edges.GroupBudget.Edges.Group)
	}
	var r *uuid.UUID
	if transaction.Edges.Request != nil {
		r = &transaction.Edges.Request.ID
	}
	return &TransactionResponse{
		ID:        transaction.ID,
		Amount:    transaction.Edges.Detail.Amount,
		Target:    transaction.Edges.Detail.Target,
		Request:   r,
		Tags:      tags,
		Group:     g,
		CreatedAt: transaction.CreatedAt,
	}
}
