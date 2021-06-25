// Code generated by entc, DO NOT EDIT.

package tag

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the tag type in the database.
	Label = "tag"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// EdgeRequest holds the string denoting the request edge name in mutations.
	EdgeRequest = "request"
	// EdgeTransaction holds the string denoting the transaction edge name in mutations.
	EdgeTransaction = "transaction"
	// Table holds the table name of the tag in the database.
	Table = "tags"
	// RequestTable is the table the holds the request relation/edge. The primary key declared below.
	RequestTable = "request_tag"
	// RequestInverseTable is the table name for the Request entity.
	// It exists in this package in order to avoid circular dependency with the "request" package.
	RequestInverseTable = "requests"
	// TransactionTable is the table the holds the transaction relation/edge. The primary key declared below.
	TransactionTable = "transaction_tag"
	// TransactionInverseTable is the table name for the Transaction entity.
	// It exists in this package in order to avoid circular dependency with the "transaction" package.
	TransactionInverseTable = "transactions"
)

// Columns holds all SQL columns for tag fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
}

var (
	// RequestPrimaryKey and RequestColumn2 are the table columns denoting the
	// primary key for the request relation (M2M).
	RequestPrimaryKey = []string{"request_id", "tag_id"}
	// TransactionPrimaryKey and TransactionColumn2 are the table columns denoting the
	// primary key for the transaction relation (M2M).
	TransactionPrimaryKey = []string{"transaction_id", "tag_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
