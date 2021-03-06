package fields

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

type JSON map[string]interface{}

// Get returns the value as a string
func (j JSON) Get(key string) string {
	return fmt.Sprintf("%v", j[key])
}

// Int64 returns the value as an int64
func (j JSON) Int64(key string) (int64, error) {
	return strconv.ParseInt(j.Get(key), 10, 64)
}

// Scan converts an SQL value into JSON
func (j *JSON) Scan(value interface{}) error {
	// Parse the bytes as JSON
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("JSON scan returned non-bytes")
	}
	return json.Unmarshal(b, j)
}

// Value returns the JSON formatted for insert into SQL
func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}
