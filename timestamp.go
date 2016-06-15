package fields

import (
	"time"

	"github.com/aodin/sol"
	"github.com/aodin/sol/postgres"
	"github.com/lib/pq"
)

// Timestamp records create and update timestamps.
type Timestamp struct {
	CreatedAt time.Time   `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at" json:"updated_at,omitempty"`
}

// Age returns the duration since the timestamp was created.
func (ts Timestamp) Age() time.Duration {
	return ts.age(time.Now().UTC())
}

func (ts Timestamp) age(now time.Time) time.Duration {
	return now.Sub(ts.CreatedAt)
}

// LastActivity returns the time of the lastest activity on the timestamp -
// either when it was updated or created. Updated is assumed to have always
// been at or after creation.
func (ts Timestamp) LastActivity() time.Time {
	if ts.WasUpdated() {
		return ts.UpdatedAt.Time
	}
	return ts.CreatedAt
}

// SetUpdatedAt sets the updated_at field
func (ts *Timestamp) SetUpdatedAt(when time.Time) {
	ts.UpdatedAt.Valid = true
	ts.UpdatedAt.Time = when
}

// Updated returns true if the timestamp has been updated
func (ts Timestamp) WasUpdated() bool {
	return ts.UpdatedAt.Valid && !ts.UpdatedAt.Time.IsZero()
}

var _ sol.Modifier = Timestamp{}

func (ts Timestamp) Modify(table sol.Tabular) error {
	// TODO Determine the column names from the struct's db tags
	columns := []sol.ColumnElem{
		sol.Column(
			"created_at",
			postgres.Timestamp().WithTimezone().NotNull().Default(postgres.Now),
		),
		sol.Column(
			"updated_at",
			postgres.Timestamp().WithTimezone(),
		),
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

// TimestampColumns returns a Modifier suitable for inclusion in a Table
func TimestampColumns() Timestamp {
	return Timestamp{}
}
