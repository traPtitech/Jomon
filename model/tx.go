package model

import (
	"context"

	"github.com/traPtitech/Jomon/ent"
)

func Tx(ctx context.Context, cli *ent.Client) (*ent.Tx, error) {
	tx, err := cli.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// RollbackWithError is a helper function to rollback a transaction and return an error.
func RollbackWithError(ctx context.Context, tx *ent.Tx, err error) error {
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Rollback()
}
