package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// GroupOwner holds the schema definition for the GroupOwner entity.
type GroupOwner struct {
	ent.Schema
}

// Fields of the GroupOwner.
func (GroupOwner) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("owner"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the GroupOwner.
func (GroupOwner) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).
			Ref("owner").
			Unique().
			Required(),
	}
}
