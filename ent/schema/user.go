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
		field.String("name").
			Unique().
			NotEmpty(),
		field.String("display_name"),
		field.Bool("account_manager").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group_user", Group.Type).
			Ref("user"),
		edge.From("group_owner", Group.Type).
			Ref("owner"),
		edge.From("comment", Comment.Type).
			Ref("user"),
		edge.From("application_status", ApplicationStatus.Type).
			Ref("user"),
		edge.From("application", Application.Type).
			Ref("user"),
		edge.From("file", File.Type).
			Ref("user"),
		edge.From("application_target", ApplicationTarget.Type).
			Ref("user"),
	}
}
