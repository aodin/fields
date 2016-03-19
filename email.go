package fields

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// Email is an email type
type Email string

// Scan converts an SQL value into an Email
func (email *Email) Scan(value interface{}) error {
	*email = Email(string(value.([]byte)))
	return nil
}

// Value returns the email as a string
func (email Email) Value() (driver.Value, error) {
	return string(email), nil
}

// UnmarshalJSON for emails trims spaces
func (email *Email) UnmarshalJSON(text []byte) error {
	b := bytes.NewBuffer(text)
	dec := json.NewDecoder(b)
	var n string
	if err := dec.Decode(&n); err != nil {
		return err
	}
	*email = Email(strings.TrimSpace(n))
	return nil
}

// Normalize will perform an in-place normalization of the email, only
// returning an email if normalization fails
func (email *Email) Normalize() error {
	normalized, err := NormalizeEmail(string(*email))
	if err != nil {
		return err
	}
	*email = Email(normalized)
	return nil
}

// NormalizeEmail will check that the email has at least one @ - it will then
// lowercase both the local name and domain
func NormalizeEmail(email string) (string, error) {
	parts := strings.Split(strings.TrimSpace(email), "@")
	if len(parts) <= 1 {
		return "", fmt.Errorf("Emails must contain a '@'")
	}
	if len(parts) > 2 {
		return "", fmt.Errorf("Emails cannot have more than one '@'")
	}
	if parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("Emails must be of the form 'user@domain'")
	}
	return strings.ToLower(strings.Join(parts, "@")), nil
}

// NewEmail creates a new Email
func NewEmail(email string) Email {
	out, _ := NormalizeEmail(email)
	return Email(out)
}
