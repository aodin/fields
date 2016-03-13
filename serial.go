package fields

import (
	sql "github.com/aodin/sol"
	pg "github.com/aodin/sol/postgres"
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

var _ sql.Modifier = Serial{}

// Modify implements the sol.Modifier interface
func (serial Serial) Modify(table *sql.TableElem) error {
	// TODO determine column name and settings from the db tags?
	return sql.Column("id", pg.Serial()).Modify(table)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
// func (serial Serial) UnmarshalText(text []byte) error {
// 	return fmt.Errorf("Serial IDs must be set by the server")
// }

// NewSerial creates a new Serial
func NewSerial(id uint64) Serial {
	return Serial{ID: id}
}
