package model

import (
	"context"
	"fmt"

	"github.com/traPtitech/Jomon/ent"
)

func setupTx(client *ent.Client, ctx context.Context) (*ent.Tx, *ent.Client, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	return tx, tx.Client(), nil
}

func rollBackWithErr(tx *ent.Tx, ctx context.Context, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%v, rolling back transaction: %v", err, rerr)
	}
	return err
}
