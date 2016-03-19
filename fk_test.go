package fields

import (
	"encoding/json"
	"testing"
)

func TestFK(t *testing.T) {
	test := struct {
		FK ImmutableFK `json:"fk_id" xml:"FKID"`
	}{}

	// Non-integer values should error
	if err := json.Unmarshal([]byte(`{"fk_id":"a"}`), &test); err == nil {
		t.Errorf("Unmarshal JSON should error when a non-number")
	}
	if test.FK.ID != 0 {
		t.Errorf("test.FK should still be zero")
	}

	if err := json.Unmarshal([]byte(`{"fk_id":0}`), &test); err == nil {
		t.Errorf("Unmarshal JSON should error when given zero")
	}
	if test.FK.ID != 0 {
		t.Errorf("test.FK should still be zero")
	}

	if err := json.Unmarshal([]byte(`{"fk_id":-1}`), &test); err == nil {
		t.Errorf(
			"Unmarshal JSON should error when given a number less than zero",
		)
	}
	if test.FK.ID != 0 {
		t.Errorf("test.FK should still be zero")
	}

	if err := json.Unmarshal([]byte(`{"fk_id":1}`), &test); err != nil {
		t.Errorf("Unmarshal JSON should not error: %s", err)
	}
	if test.FK.ID != 1 {
		t.Errorf("unexpected test.FK: %d != 1", test.FK.ID)
	}

	if err := json.Unmarshal([]byte(`{"fk_id":1}`), &test); err != nil {
		t.Errorf("Unmarshal JSON should not error if being set with the same value: %s", err)
	}
	if test.FK.ID != 1 {
		t.Errorf("unexpected test.FK: %d != 1", test.FK.ID)
	}

	if err := json.Unmarshal([]byte(`{"fk_id":2}`), &test); err == nil {
		t.Errorf(
			"Unmarshal JSON should error if being set with a different value",
		)
	}
	if test.FK.ID != 1 {
		t.Errorf("unexpected test.FK: %d != 1", test.FK.ID)
	}
}
