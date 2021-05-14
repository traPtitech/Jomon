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
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the RequestTag.
func (RequestTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("request", Request.Type).
			Ref("tag").
			Unique().
			Required(),
		edge.From("tag", Tag.Type).
			Ref("request_tag").
			Unique().
			Required(),
	}
}
