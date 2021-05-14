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
		field.String("target").
			Default(""),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the TransactionDetail.
func (TransactionDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("detail").
			Unique().
			Required(),
		edge.From("request", Request.Type).
			Ref("transaction_detail").
			Unique(),
		//edge.From("group", Group.Type).
		//	Ref("tag").
		//	Unique(),
	}
}
