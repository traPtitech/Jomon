package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RequestTarget holds the schema definition for the RequestTarget entity.
type RequestTarget struct {
	ent.Schema
}

// Fields of the RequestTarget.
func (RequestTarget) Fields() []ent.Field {
	return []ent.Field{
		field.String("target"),
		field.Int("request_id"),
		field.Time("paid_at").
			Nillable().
			Optional(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the RequestTarget.
func (RequestTarget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("request", Request.Type).
			Field("request_id").
			Unique().
			Required(),
	}
}
