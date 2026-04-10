package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// UserSubject holds the schema definition for the UserSubject entity.
type UserSubject struct {
	ent.Schema
}

// Fields of the UserSubject.
func (UserSubject) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			StorageKey("subject").
			MaxLen(36),
		field.UUID("user_id", uuid.UUID{}),
	}
}

// Edges of the UserSubject.
func (UserSubject) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Field("user_id").
			Unique().
			Required(),
	}
}
