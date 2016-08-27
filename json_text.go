package pq_types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONText is a raw encoded JSON value, compatible with PostgreSQL's varchar, text, json and jsonb.
// It behaves like json.RawMessage by implementing json.Marshaler and json.Unmarshaler
// and can be used to delay JSON decoding or precompute a JSON encoding.
type JSONText []byte

// String implements fmt.Stringer for better output and logging.
func (j JSONText) String() string {
	return string(j)
}

// MarshalJSON returns j as the JSON encoding of j.
func (j JSONText) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte(`null`), nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data.
func (j *JSONText) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONText.UnmarshalJSON: on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil

}

// Value implements database/sql/driver Valuer interface.
// It performs basic validation by unmarshaling itself into json.RawMessage.
// If j is not valid JSON, it returns and error.
func (j JSONText) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	var m json.RawMessage
	var err = json.Unmarshal(j, &m)
	if err != nil {
		return []byte{}, err
	}
	return []byte(j), nil
}

// Scan implements database/sql Scanner interface.
// It store value in *j. No validation is done.
func (j *JSONText) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("JSONText.Scan: expected []byte or string, got %T (%q)", value, value)
	}

	*j = JSONText(append((*j)[0:0], b...))
	return nil
}

// check interfaces
var (
	_ json.Marshaler   = JSONText{}
	_ json.Unmarshaler = &JSONText{}
	_ driver.Valuer    = JSONText{}
	_ sql.Scanner      = &JSONText{}
	_ fmt.Stringer     = JSONText{}
	_ fmt.Stringer     = &JSONText{}
)
