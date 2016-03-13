package fields

import (
	"testing"

	sql "github.com/aodin/sol"
)

type serialTest struct {
	Serial
	Name string
}

var SerialTests = sql.Table("serial_tests",
	Serial{},
)

func TestSerial(t *testing.T) {
	var item serialTest

	if item.Exists() {
		t.Errorf("Item should not exist until ID is set")
	}

	item.ID = 1
	if !item.Exists() {
		t.Errorf("Item should exist once ID is set")
	}
}
