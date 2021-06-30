// Code generated by entc, DO NOT EDIT.

package transactiondetail

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the transactiondetail type in the database.
	Label = "transaction_detail"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAmount holds the string denoting the amount field in the database.
	FieldAmount = "amount"
	// FieldTarget holds the string denoting the target field in the database.
	FieldTarget = "target"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeTransaction holds the string denoting the transaction edge name in mutations.
	EdgeTransaction = "transaction"
	// Table holds the table name of the transactiondetail in the database.
	Table = "transaction_details"
	// TransactionTable is the table the holds the transaction relation/edge.
	TransactionTable = "transaction_details"
	// TransactionInverseTable is the table name for the Transaction entity.
	// It exists in this package in order to avoid circular dependency with the "transaction" package.
	TransactionInverseTable = "transactions"
	// TransactionColumn is the table column denoting the transaction relation/edge.
	TransactionColumn = "transaction_detail"
)

// Columns holds all SQL columns for transactiondetail fields.
var Columns = []string{
	FieldID,
	FieldAmount,
	FieldTarget,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "transaction_details"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"transaction_detail",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultAmount holds the default value on creation for the "amount" field.
	DefaultAmount int
	// DefaultTarget holds the default value on creation for the "target" field.
	DefaultTarget string
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)