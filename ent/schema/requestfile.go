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
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the RequestFile.
func (RequestFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("request", Request.Type).
			Ref("file").
			Unique().
			Required(),
		edge.To("file", File.Type).
			Unique(),
	}
}
