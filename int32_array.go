package pq_types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Int32Array is a slice of int32 values, compatible with PostgreSQL's int[] and intarray module.
type Int32Array []int32

func (a Int32Array) Len() int           { return len(a) }
func (a Int32Array) Less(i, j int) bool { return a[i] < a[j] }
func (a Int32Array) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Value implements database/sql/driver Valuer interface.
func (a Int32Array) Value() (driver.Value, error) {
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
func (a *Int32Array) Scan(value interface{}) error {
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
		return fmt.Errorf("Int32Array.Scan: expected []byte or string, got %T (%q)", value, value)
	}

	if len(b) < 2 || b[0] != '{' || b[len(b)-1] != '}' {
		return fmt.Errorf("Int32Array.Scan: unexpected data %q", b)
	}

	p := strings.Split(string(b[1:len(b)-1]), ",")

	// reuse underlying array if present
	if *a == nil {
		*a = make(Int32Array, 0, len(p))
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
		*a = append(*a, int32(i))
	}

	return nil
}

// EqualWithoutOrder returns true if two int32 arrays are equal without order, false otherwise.
// It may sort both arrays in-place to do so.
func (a Int32Array) EqualWithoutOrder(b Int32Array) bool {
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
	_ sort.Interface = Int32Array{}
	_ driver.Valuer  = Int32Array{}
	_ sql.Scanner    = &Int32Array{}
)
