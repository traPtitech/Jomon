package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Transaction holds the schema definition for the Transaction entity.
type Transaction struct {
	ent.Schema
}

// Fields of the Transaction.
func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Transaction.
func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("detail", TransactionDetail.Type).
			Unique(),
		edge.To("tag", Tag.Type),
		edge.From("group_budget", GroupBudget.Type).
			Ref("transaction").
			Unique(),
		edge.From("request", Request.Type).
			Ref("transaction").
			Unique(),
	}
}
