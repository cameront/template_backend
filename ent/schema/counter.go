package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Counter holds the schema definition for the Counter entity.
type Counter struct {
	ent.Schema
}

// Fields of the Counter.
func (Counter) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.Int64("value"),
	}
}

// Edges of the Counter.
func (Counter) Edges() []ent.Edge {
	return nil
}
