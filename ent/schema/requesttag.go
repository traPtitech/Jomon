package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RequestTag holds the schema definition for the RequestTag entity.
type RequestTag struct {
	ent.Schema
}

// Fields of the RequestTag.
func (RequestTag) Fields() []ent.Field {
	return []ent.Field{
		field.Int("request_id"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the RequestTag.
func (RequestTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("request", Request.Type).
			Field("request_id").
			Unique().
			Required(),
		edge.To("tag", Tag.Type).
			Unique(),
	}
}
