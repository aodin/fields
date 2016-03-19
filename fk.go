package fields

import (
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/aodin/sol"
	"github.com/aodin/sol/types"
)

// ImmutableFK is a foreign key. It embeds an ID and adds a column name
// (such as remote_id) and foreign key table (such as Remotes)
type ImmutableFK struct {
	ID    uint64
	Name  string
	Table *sol.TableElem
}

var _ sol.Modifier = ImmutableFK{}

// Exists returns true if the foreign key table has an entry with the given ID
func (fk ImmutableFK) Exists(conn sol.Conn) bool {
	// If the table is missing, the key must not exist
	if fk.Table == nil {
		return false
	}
	var id int64
	stmt := sol.Select(
		fk.Table.C("id"),
	).Where(
		fk.Table.C("id").Equals(fk.ID),
	).Limit(1)
	conn.Query(stmt, &id)
	return id != 0
}

// MarshalJSON returns the inner ID
func (fk ImmutableFK) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(fk.ID, 10)), nil
}

// Scan implements the database/sql.Scanner interface
func (fk *ImmutableFK) Scan(value interface{}) error {
	fk.ID = uint64(value.(int64))
	return nil
}

// Value implements the database/sql/driver.Valuer interface
func (fk ImmutableFK) Value() (driver.Value, error) {
	return int64(fk.ID), nil
}

// UnmarshalJSON allows an ID to be set once, but only once
// TODO Why not use UnmarshalText?
// https://github.com/golang/go/issues/9650
func (fk *ImmutableFK) UnmarshalJSON(b []byte) error {
	id, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse foreign key: %s", err)
	}
	if fk.ID != 0 {
		if fk.ID == id {
			// Foreign key is unchanged - this is okay
			return nil
		}
		return fmt.Errorf("foreign keys cannot be overwritten")
	}
	if id <= 0 {
		return fmt.Errorf("foreign keys cannot be zero or less than zero")
	}
	fk.ID = id
	return err
}

// Modify implements the sol.Modifier interface
func (fk ImmutableFK) Modify(table *sol.TableElem) error {
	return sol.ForeignKey(
		fk.Name,
		fk.Table.C("id"),
		types.Integer().NotNull(),
	).OnUpdate(sol.Cascade).OnDelete(sol.Cascade).Modify(table)
}

// SetTable sets the column name and table of the foreign key
func (fk *ImmutableFK) SetTable(name string, table *sol.TableElem) {
	fk.Name = name
	fk.Table = table
	return
}
