package pq_types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"
)

// StringArray is a slice of string values, compatible with PostgreSQL's varchar[].
type StringArray []string

func (a StringArray) Len() int           { return len(a) }
func (a StringArray) Less(i, j int) bool { return a[i] < a[j] }
func (a StringArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Value implements database/sql/driver Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	res := make([]string, len(a))
	for i, e := range a {
		r := e
		r = strings.Replace(r, `\`, `\\`, -1)
		r = strings.Replace(r, `"`, `\"`, -1)
		res[i] = `"` + r + `"`
	}
	return []byte("{" + strings.Join(res, ",") + "}"), nil
}

// Scan implements database/sql Scanner interface.
func (a *StringArray) Scan(value interface{}) error {
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
		return fmt.Errorf("StringArray.Scan: expected []byte or string, got %T (%q)", value, value)
	}

	if len(b) < 2 || b[0] != '{' || b[len(b)-1] != '}' {
		return fmt.Errorf("StringArray.Scan: unexpected data %q", b)
	}

	// reuse underlying array if present
	if *a == nil {
		*a = make(StringArray, 0)
	}
	*a = (*a)[:0]

	if len(b) == 2 { // '{}'
		return nil
	}

	reader := bytes.NewReader(b[1 : len(b)-1]) // skip '{' and '}'

	// helper function to read next rune and check if it valid
	readRune := func() (rune, error) {
		r, _, err := reader.ReadRune()
		if err != nil {
			return 0, err
		}
		if r == unicode.ReplacementChar {
			return 0, fmt.Errorf("StringArray.Scan: invalid rune")
		}
		return r, nil
	}

	var q bool
	var e []rune
	for {
		// read next rune and check if we are done
		r, err := readRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch r {
		case '"':
			// enter or leave quotes
			q = !q
			continue
		case ',':
			// end of element unless in we are in quotes
			if !q {
				*a = append(*a, string(e))
				e = e[:0]
				continue
			}
		case '\\':
			// skip to next rune, it should be present
			n, err := readRune()
			if err != nil {
				return err
			}
			r = n
		}

		e = append(e, r)
	}

	// we should not be in quotes at this point
	if q {
		panic("StringArray.Scan bug")
	}

	// add last element
	*a = append(*a, string(e))
	return nil
}

// check interfaces
var (
	_ sort.Interface = StringArray{}
	_ driver.Valuer  = StringArray{}
	_ sql.Scanner    = &StringArray{}
)
