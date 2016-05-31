package fields

import (
	"github.com/aodin/sol"
	"github.com/aodin/sol/postgres"
)

// Serial is an embeddable type that can be used as a primary key for a table
type Serial struct {
	ID uint64 `db:"id,omitempty" json:"id" xml:"ID"`
}

// Exists returns true if the Pk is non-zero
func (serial Serial) Exists() bool {
	return serial.ID != 0
}

// GetID returns the ID as a uint64
func (serial Serial) GetID() uint64 {
	return serial.ID
}

func (serial Serial) Keys() []interface{} {
	return []interface{}{serial.ID}
}

var _ sol.Modifier = Serial{}

// Modify implements the sol.Modifier interface
func (serial Serial) Modify(table sol.Tabular) error {
	// TODO determine column name and settings from the db tags?
	return sol.Column("id", postgres.Serial()).Modify(table)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
// func (serial Serial) UnmarshalText(text []byte) error {
// 	return fmt.Errorf("Serial IDs must be set by the server")
// }

// NewSerial creates a new Serial
func NewSerial(id uint64) Serial {
	return Serial{ID: id}
}
