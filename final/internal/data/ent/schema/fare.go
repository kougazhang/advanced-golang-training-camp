package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// Fare holds the schema definition for the Fare entity.
type Fare struct {
	ent.Schema
}

// Fields of the Fare.
func (Fare) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("org_airport"),
		field.String("arr_airport"),
		field.String("passage_type"),
		field.Time("first_travel_date").
			Default(time.Now).SchemaType(map[string]string{
			dialect.SQLite: "datetime",
			dialect.MySQL:  "datetime",
		}),
		field.Time("last_travel_date").
			Default(time.Now).SchemaType(map[string]string{
			dialect.SQLite: "datetime",
			dialect.MySQL:  "datetime",
		}),
		field.Float("amount").
			GoType(float64(0)).
			SchemaType(map[string]string{
				dialect.SQLite: "double",
				dialect.MySQL:  "decimal(6,2)",
			}),
		field.Time("created_at").
			Default(time.Now).SchemaType(map[string]string{
			dialect.SQLite: "datetime",
			dialect.MySQL:  "datetime",
		}),
		field.Time("updated_at").
			Default(time.Now).SchemaType(map[string]string{
			dialect.SQLite: "datetime",
			dialect.MySQL:  "datetime",
		}),
	}
}

// Edges of the Fare.
func (Fare) Edges() []ent.Edge {
	return nil
}
