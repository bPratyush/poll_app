package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

// Fields of the Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.String("message").
			NotEmpty(),
		field.String("type").
			Default("vote_changed"), // vote_changed, poll_updated, etc.
		field.Int("poll_id").
			Optional(),
		field.Bool("read").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Notification.
func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("notifications").
			Unique().
			Required(),
	}
}
