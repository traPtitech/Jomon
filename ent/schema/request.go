package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Request holds the schema definition for the Request entity.
type Request struct {
	ent.Schema
}

// Fields of the Request.
func (Request) Fields() []ent.Field {
	return []ent.Field{
		field.String("created_by"),
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
		edge.To("file", RequestFile.Type),
		edge.To("tag", RequestTag.Type),
		edge.To("transaction_detail", TransactionDetail.Type),
		edge.To("comment", Comment.Type),
	}
}
