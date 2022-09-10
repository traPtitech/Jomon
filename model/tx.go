package model

import (
	"github.com/traPtitech/Jomon/ent"
)

// RollbackWithError is a helper function to rollback a transaction and return an error.
func RollbackWithError(tx *ent.Tx, err error) error {
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Rollback()
}
