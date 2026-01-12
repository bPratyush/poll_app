package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Vote holds the schema definition for the Vote entity.
type Vote struct {
	ent.Schema
}

// Fields of the Vote.
func (Vote) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Vote.
func (Vote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("votes").
			Unique().
			Required(),
		edge.From("option", PollOption.Type).
			Ref("votes").
			Unique().
			Required(),
	}
}

// Indexes of the Vote.
func (Vote) Indexes() []ent.Index {
	return []ent.Index{
		// Ensure a user can only vote once per poll option
		index.Edges("user", "option").
			Unique(),
	}
}
