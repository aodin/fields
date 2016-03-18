package fields

import (
	"testing"
)

func TestJSON(t *testing.T) {
	// TODO test invalid JSON types
	j := JSON{
		"int":     1,
		"float":   2.1,
		"string":  "3",
		"boolean": true,
		"nil":     nil,
	}
	if j.Get("int") != "1" {
		t.Errorf("unexpected int value: %s != 1", j.Get("int"))
	}
	if j.Get("float") != "2.1" {
		t.Errorf("unexpected float value: %s != 2.1", j.Get("float"))
	}
	if j.Get("string") != "3" {
		t.Errorf("unexpected string value: %s != 3", j.Get("string"))
	}
	if j.Get("boolean") != "true" {
		t.Errorf("unexpected boolean value: %s != true", j.Get("boolean"))
	}
	if j.Get("nil") != "<nil>" {
		t.Errorf("unexpected nil value: %s != <nil>", j.Get("nil"))
	}
}
