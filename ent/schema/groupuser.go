package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// GroupUser holds the schema definition for the GroupUser entity.
type GroupUser struct {
	ent.Schema
}

// Fields of the GroupUser.
func (GroupUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the GroupUser.
func (GroupUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).
			Ref("user").
			Unique().
			Required(),
	}
}
