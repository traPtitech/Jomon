package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RequestFile holds the schema definition for the RequestFile entity.
type RequestFile struct {
	ent.Schema
}

// Fields of the RequestFile.
func (RequestFile) Fields() []ent.Field {
	return []ent.Field{
		field.Int("request_id"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the RequestFile.
func (RequestFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("request", Request.Type).
			Field("request_id").
			Unique().
			Required(),
		edge.To("file", File.Type).
			Unique().
			Required(),
	}
}
