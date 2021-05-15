package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// RequestStatus holds the schema definition for the RequestStatus entity.
type RequestStatus struct {
	ent.Schema
}

// Fields of the RequestStatus.
func (RequestStatus) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
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
		edge.From("request", Request.Type).
			Ref("status").
			Unique().
			Required(),
		edge.To("user", User.Type).
			Unique().
			Required(),
	}
}
