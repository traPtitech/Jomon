package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TransactionTag holds the schema definition for the TransactionTag entity.
type TransactionTag struct {
	ent.Schema
}

// Fields of the TransactionTag.
func (TransactionTag) Fields() []ent.Field {
	return []ent.Field{
		field.Int("transaction_id"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the TransactionTag.
func (TransactionTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("transaction", Transaction.Type).
			Field("transaction_id").
			Unique().
			Required(),
		edge.To("tag", Tag.Type).
			Unique(),
	}
}
