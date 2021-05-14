package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Administrator holds the schema definition for the Administrator entity.
type Administrator struct {
	ent.Schema
}

// Fields of the Administrator.
func (Administrator) Fields() []ent.Field {
	return []ent.Field{
		field.String("trap_id"),
	}
}

// Edges of the Administrator.
func (Administrator) Edges() []ent.Edge {
	return nil
}
