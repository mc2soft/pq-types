package pq_types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Int64Array is a slice of int32 values, compatible with PostgreSQL's int[] and intarray module.
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

	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Int64Array.Scan: expected []byte, got %T (%q)", value, value)
	}

	if len(v) < 2 || v[0] != '{' || v[len(v)-1] != '}' {
		return fmt.Errorf("Int64Array.Scan: unexpected data %q", v)
	}

	p := strings.Split(string(v[1:len(v)-1]), ",")
	res := make(Int64Array, 0, len(p))
	for _, s := range p {
		if s == "" {
			continue
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		res = append(res, int64(i))
	}
	*a = res
	return nil
}

// EqualWithoutOrder returns true if two int32 arrays are equal without order, false otherwise.
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
