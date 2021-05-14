package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RequestStatus holds the schema definition for the RequestStatus entity.
type RequestStatus struct {
	ent.Schema
}

// Fields of the RequestStatus.
func (RequestStatus) Fields() []ent.Field {
	return []ent.Field{
		field.String("created_by"),
		field.Int("request_id"),
		field.Enum("status").
			Values("submitted", "fix_required", "accepted", "completed", "rejected").
			Default("submitted"),
		field.String("reason"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the RequestStatus.
func (RequestStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("request", Request.Type).
			Field("request_id").
			Unique().
			Required(),
	}
}
