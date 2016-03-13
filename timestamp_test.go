package fields

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	now := time.Date(2015, 3, 2, 0, 0, 0, 0, time.UTC)
	then := time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC)
	again := time.Date(2015, 3, 3, 0, 0, 0, 0, time.UTC)

	ts := newTimestamp(then)
	if !ts.LastActivity().Equal(then) {
		t.Errorf("Last activity should be equal to its creation")
	}
	if ts.IsDeleted() {
		t.Errorf("Timestamp should not be deleted")
	}
	if ts.WasUpdated() {
		t.Errorf("Timestamp should not be updated")
	}

	if (24 * time.Hour) != ts.age(now) {
		t.Errorf("Timestamp should be one day old")
	}

	ts.UpdatedAt = &again
	if !ts.LastActivity().Equal(again) {
		t.Errorf("Last activity should be equal to its update")
	}

	ts.setDeletedAt(now)
	if !ts.LastActivity().Equal(now) {
		t.Errorf("Last activity should be equal to now")
	}

	if !ts.IsDeleted() {
		t.Errorf("Timestamp should be deleted")
	}
}
