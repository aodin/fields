package fields

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid := NewUUID()

	// Generated UUIDs should be valid
	parsed, err := ParseUUID(uuid.String())
	if err != nil {
		t.Fatalf("ParseUUID should not error with a valid UUID: %s", err)
	}
	if !uuid.Equals(parsed) {
		t.Errorf("UUID should equal the parsed UUID")
	}
	if !uuid.Exists() {
		t.Errorf("UUID should exist")
	}

	b, err := json.Marshal(uuid)
	if err != nil {
		t.Fatalf("JSON marshal should not error: %s", err)
	}
	out := fmt.Sprintf(`"%s"`, uuid.String())
	if out != string(b) {
		t.Errorf("Execpected JSON marshal output: %s != %s", out, string(b))
	}
}
