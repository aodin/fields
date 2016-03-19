package fields

import "testing"

func TestEmail(t *testing.T) {
	email := NewEmail("A@example.com")
	if string(email) != "a@example.com" {
		t.Errorf("unexpected email: %s != a@example.com", email)
	}

	tests := []struct {
		in, out string
	}{
		{in: "", out: ""},
		{in: "dachshundlover", out: ""},
		{in: "k@r@j", out: ""},
		{in: "K@J", out: "k@j"},
		{in: "K@", out: ""},
	}
	for _, test := range tests {
		result, err := NormalizeEmail(test.in)
		if result != test.out {
			t.Errorf(
				"unexpected normalized email %s != %s",
				result, test.out,
			)
			continue
		}
		if result == "" && err == nil {
			t.Errorf("Email is invalid but no error is present")
		}
	}
}

func TestEmail_Normalize(t *testing.T) {
	// Normalize the email in place
	email := Email("A@Example.com")
	if err := email.Normalize(); err != nil {
		t.Fatalf("email normalization should not fail: %s", err)
	}
	if string(email) != "a@example.com" {
		t.Errorf("unexpected email: %s != a@example.com", email)
	}
}
