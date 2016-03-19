package fields

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// URL is a URL field
// TODO Just boilerplate for now
type URL string

// Scan converts an SQL value into a URL
func (url *URL) Scan(value interface{}) error {
	*url = URL(string(value.([]byte)))
	return nil
}

// Value returns the URL as a string
func (url URL) Value() (driver.Value, error) {
	return string(url), nil
}

// UnmarshalJSON
func (url *URL) UnmarshalJSON(text []byte) error {
	b := bytes.NewBuffer(text)
	dec := json.NewDecoder(b)
	var n string
	if err := dec.Decode(&n); err != nil {
		return err
	}
	*url = URL(strings.TrimSpace(n))
	return nil
}

// NewURL creates a new URL
func NewURL(url string) URL {
	return URL(url)
}
