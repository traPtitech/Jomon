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
		field.Time("created_at").
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
		edge.To("transaction_detail", TransactionDetail.Type),
		edge.To("comment", Comment.Type),
		edge.To("user", User.Type).
			Unique(),
	}
}
