package pq_types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Int64Array is a slice of int64 values, compatible with PostgreSQL's bigint[].
type Int64Array []int64

func (a Int64Array) Len() int           { return len(a) }
func (a Int64Array) Less(i, j int) bool { return a[i] < a[j] }
func (a Int64Array) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Value implements database/sql/driver Valuer interface.
func (a Int64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	s := make([]string, len(a))
	for i, v := range a {
		s[i] = strconv.Itoa(int(v))
	}
	return []byte("{" + strings.Join(s, ",") + "}"), nil
}

// Scan implements database/sql Scanner interface.
func (a *Int64Array) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("Int64Array.Scan: expected []byte or string, got %T (%q)", value, value)
	}

	if len(b) < 2 || b[0] != '{' || b[len(b)-1] != '}' {
		return fmt.Errorf("Int64Array.Scan: unexpected data %q", b)
	}

	p := strings.Split(string(b[1:len(b)-1]), ",")

	// reuse underlying array if present
	if *a == nil {
		*a = make(Int64Array, 0, len(p))
	}
	*a = (*a)[:0]

	for _, s := range p {
		if s == "" {
			continue
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*a = append(*a, int64(i))
	}

	return nil
}

// EqualWithoutOrder returns true if two int64 arrays are equal without order, false otherwise.
// It may sort both arrays in-place to do so.
func (a Int64Array) EqualWithoutOrder(b Int64Array) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Sort(a)
	sort.Sort(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// check interfaces
var (
	_ sort.Interface = Int64Array{}
	_ driver.Valuer  = Int64Array{}
	_ sql.Scanner    = &Int64Array{}
)
