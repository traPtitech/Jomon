package model

import "github.com/traPtitech/Jomon/ent"

func ConvertEntTransactionToModelTransaction(transaction *ent.Transaction) *Transaction {
	return &Transaction{
		ID:        transaction.ID,
		CreatedAt: transaction.CreatedAt,
	}
}
