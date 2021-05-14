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
		edge.From("status", RequestStatus.Type).
			Ref("request"),
		edge.From("target", RequestTarget.Type).
			Ref("request"),
		edge.From("file", RequestFile.Type).
			Ref("request"),
		edge.From("tag", RequestTag.Type).
			Ref("request"),
		edge.From("transaction_detail", TransactionDetail.Type).
			Ref("request"),
		edge.From("comment", Comment.Type).
			Ref("request"),
	}
}
