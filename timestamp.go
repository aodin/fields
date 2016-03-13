package fields

import (
	"time"

	sql "github.com/aodin/sol"
	pg "github.com/aodin/sol/postgres"
	"github.com/aodin/sol/types"
)

// Timestamp records creation, update, and optional deletion timestamps.
type Timestamp struct {
	CreatedAt time.Time  `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// IsDeleted returns true if the object has been deleted
func (ts Timestamp) IsDeleted() bool {
	return ts.DeletedAt != nil
}

// SetDeletedAt modifies the timestamp struct to set DeletedAt to now
func (ts *Timestamp) SetDeletedAt() time.Time {
	return ts.setDeletedAt(time.Now().UTC())
}

func (ts *Timestamp) setDeletedAt(t time.Time) time.Time {
	ts.DeletedAt = &t
	return t
}

// Updated returns true if the timestamp has been updated
func (ts Timestamp) WasUpdated() bool {
	return ts.UpdatedAt != nil
}

// Age returns the duration since the timestamp was created.
func (ts Timestamp) Age() time.Duration {
	return ts.age(time.Now().UTC())
}

func (ts Timestamp) age(now time.Time) time.Duration {
	return now.Sub(ts.CreatedAt)
}

// LastActivity returns the time of the lastest activity on the timestamp -
// either when it was last deleted, updated, or created
func (ts Timestamp) LastActivity() time.Time {
	if ts.DeletedAt != nil {
		return *ts.DeletedAt
	}
	if ts.UpdatedAt != nil {
		return *ts.UpdatedAt
	}
	return ts.CreatedAt
}

var _ sql.Modifier = Timestamp{}

func (ts Timestamp) Modify(table *sql.TableElem) error {
	// TODO Determine the column names from the struct's db tags
	columns := []sql.ColumnElem{
		sql.Column("created_at", pg.Timestamp().NotNull().Default(pg.Now)),
		sql.Column("updated_at", types.Timestamp()),
		sql.Column("deleted_at", types.Timestamp()),
	}
	for _, column := range columns {
		if err := column.Modify(table); err != nil {
			return err
		}
	}
	return nil
}

// Timestamps should only be created by the database. This constructor should
// only be used for testing.
func newTimestamp(now time.Time) Timestamp {
	return Timestamp{CreatedAt: now}
}
