package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// RequestTarget holds the schema definition for the RequestTarget entity.
type RequestTarget struct {
	ent.Schema
}

// Fields of the RequestTarget.
func (RequestTarget) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("amount"),
		field.Time("paid_at").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the RequestTarget.
func (RequestTarget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("request", Request.Type).
			Ref("target").
			Unique().
			Required(),
		edge.To("user", User.Type).
			Unique().
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.NoAction,
			}),
	}
}
