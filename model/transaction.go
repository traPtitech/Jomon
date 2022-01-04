//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	GetTransactions(ctx context.Context, query TransactionQuery) ([]*TransactionResponse, error)
	CreateTransaction(ctx context.Context, Amount int, Target string, tags []*uuid.UUID, group *uuid.UUID) (*TransactionResponse, error)
	GetTransaction(ctx context.Context, transactionID uuid.UUID) (*TransactionResponse, error)
	UpdateTransaction(ctx context.Context, transactionID uuid.UUID, Amount int, Target string, tags []*uuid.UUID, group *uuid.UUID) (*TransactionResponse, error)
}

type Transaction struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

type TransactionResponse struct {
	ID        uuid.UUID
	Amount    int
	Target    string
	Tags      []*Tag
	Group     *Group
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionQuery struct {
	Sort   *string
	Target *string
	Year   *int
	Since  *time.Time
	Until  *time.Time
	Tag    *string
	Group  *string
}
