package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ApplicationTarget holds the schema definition for the ApplicationTarget entity.
type ApplicationTarget struct {
	ent.Schema
}

// Fields of the ApplicationTarget.
func (ApplicationTarget) Fields() []ent.Field {
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

// Edges of the ApplicationTarget.
func (ApplicationTarget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("application", Application.Type).
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
