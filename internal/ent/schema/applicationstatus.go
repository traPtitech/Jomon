package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ApplicationStatus holds the schema definition for the ApplicationStatus entity.
type ApplicationStatus struct {
	ent.Schema
}

// Fields of the ApplicationStatus.
func (ApplicationStatus) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Enum("status").
			Values("submitted", "fix_required", "accepted", "completed", "rejected").
			Default("submitted"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the ApplicationStatus.
func (ApplicationStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("application", Application.Type).
			Ref("status").
			Unique().
			Required(),
		edge.To("user", User.Type).
			Unique().
			Required(),
	}
}
