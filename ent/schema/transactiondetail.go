package schema

import (
	"errors"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TransactionDetail holds the schema definition for the TransactionDetail entity.
type TransactionDetail struct {
	ent.Schema
}

// Fields of the TransactionDetail.
func (TransactionDetail) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title").NotEmpty().MaxLen(64).Validate(func(s string) error {
			if strings.Contains(s, "\n") {
				return errors.New("title cannot contain new line")
			}
			return nil
		}),
		field.Int("amount").
			Default(0),
		field.String("target").
			Default(""),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the TransactionDetail.
func (TransactionDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("detail").
			Unique(),
	}
}
