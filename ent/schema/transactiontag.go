package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// TransactionTag holds the schema definition for the TransactionTag entity.
type TransactionTag struct {
	ent.Schema
}

// Fields of the TransactionTag.
func (TransactionTag) Fields() []ent.Field {
	return nil
}

// Edges of the TransactionTag.
func (TransactionTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("tag").
			Unique().
			Required(),
		edge.From("tag", Tag.Type).
			Ref("transaction_tag").
			Unique().
			Required(),
	}
}
