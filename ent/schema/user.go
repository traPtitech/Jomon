package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("trap_id").
			Unique(),
		field.String("name"),
		field.Bool("admin").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).
			Ref("user"),
		edge.From("comment", Comment.Type).
			Ref("user").
			Unique(),
		edge.From("request_status", RequestStatus.Type).
			Ref("user").
			Unique(),
		edge.From("request", Request.Type).
			Ref("user").
			Unique(),
	}
}
