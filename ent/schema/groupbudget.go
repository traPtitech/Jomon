package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// GroupBudget holds the schema definition for the GroupBudget entity.
type GroupBudget struct {
	ent.Schema
}

// Fields of the GroupBudget.
func (GroupBudget) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("amount"),
		field.Text("comment").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the GroupBudget.
func (GroupBudget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).
			Ref("group_budget").
			Unique().
			Required(),
		edge.To("transaction", Transaction.Type),
	}
}
