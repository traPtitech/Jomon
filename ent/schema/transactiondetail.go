package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TransactionDetail holds the schema definition for the TransactionDetail entity.
type TransactionDetail struct {
	ent.Schema
}

// Fields of the TransactionDetail.
func (TransactionDetail) Fields() []ent.Field {
	return []ent.Field{
		field.Int("amount").
			Default(0),
		field.Int("request_id").
			Nillable().
			Optional(),
		field.Int("group_id").
			Nillable().
			Optional(),
		field.String("target").
			Default(""),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the TransactionDetail.
func (TransactionDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("transaction", Transaction.Type).
			Unique().
			Required(),
		edge.To("request", Request.Type).
			Field("request_id").
			Unique(),
		edge.To("group", Group.Type).
			Field("group_id").
			Unique(),
	}
}
