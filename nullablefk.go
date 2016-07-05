package fields

import (
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/aodin/sol"
	"github.com/aodin/sol/types"
)

// NullableFK is a foreign key that can be NULL. It embeds an ImmutableFK
type NullableFK struct {
	ImmutableFK
	Valid bool
}

var _ sol.Modifier = NullableFK{}

// Exists returns true if the foreign key table has an entry with the given ID
// or the FK is NULL.
func (fk NullableFK) Exists(conn sol.Conn) bool {
	if !fk.Valid {
		return true
	}
	return fk.ImmutableFK.Exists(conn)
}

// Scan converts the raw SQL value into a NullableFK
func (fk *NullableFK) Scan(value interface{}) error {
	if value == nil {
		fk.Valid = false
		return nil
	}
	if err := fk.ImmutableFK.Scan(value); err != nil {
		return err
	}
	fk.Valid = true
	return nil
}

// Value returns the FK ID or nil if the FK is not valid
func (fk NullableFK) Value() (driver.Value, error) {
	if !fk.Valid || fk.ID == 0 {
		return nil, nil
	}
	return fk.ImmutableFK.Value()
}

// MarshalJSON returns the JSON output of the FK
func (fk NullableFK) MarshalJSON() ([]byte, error) {
	if fk.Valid {
		return fk.ImmutableFK.MarshalJSON()
	}
	return []byte(`null`), nil
}

func (fk *NullableFK) UnmarshalJSON(b []byte) error {
	fk.Valid = false
	if string(b) == "null" {
		fk.ImmutableFK.ID = 0
		return nil
	}
	id, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse foreign key: %s", err)
	}
	if id < 0 {
		// Unlike ImmutableFK, zero is an acceptable value
		return fmt.Errorf("foreign keys cannot be less than zero")
	}
	fk.ImmutableFK.ID = id
	fk.Valid = id != 0
	return nil
}

// Modify implements the sol.Modifier interface
func (fk NullableFK) Modify(table sol.Tabular) error {
	return sol.ForeignKey(
		fk.Name,
		fk.Table.C("id"),
		types.Integer(),
	).OnUpdate(sol.Cascade).OnDelete(sol.Cascade).Modify(table)
}
