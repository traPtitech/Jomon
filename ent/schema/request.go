package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Request holds the schema definition for the Request entity.
type Request struct {
	ent.Schema
}

// Fields of the Request.
func (Request) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("amount"),
		field.String("title"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
	}
}

// Edges of the Request.
func (Request) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("status", RequestStatus.Type),
		edge.To("target", RequestTarget.Type),
		edge.To("file", File.Type),
		edge.To("tag", Tag.Type),
		edge.To("transaction", Transaction.Type),
		edge.To("comment", Comment.Type),
		edge.To("user", User.Type).
			Unique(),
		edge.From("group", Group.Type).
			Ref("request").
			Unique(),
	}
}
