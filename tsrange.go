package pq_types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// TimeBound represents Upper and Lower bound for TSRange.
// Time may be nil, so this will be infinity.
// If Time is nil and bound is Inclusive, it will be converted to exclusive by postgresql.
type TimeBound struct {
	Inclusive bool
	Time      *time.Time
}

// TSRange is a wrapper for postresql tsrange type.
type TSRange struct {
	LowerBound TimeBound
	UpperBound TimeBound
}

const (
	timeFormat = "2006-01-02 15:04:05"
)

// Value implements database/sql/driver Valuer interface.
func (t TSRange) Value() (driver.Value, error) {
	res := []byte{}
	if t.LowerBound.Inclusive {
		res = append(res, '[')
	} else {
		res = append(res, '(')
	}
	if t.LowerBound.Time != nil {
		tstr := t.LowerBound.Time.UTC().Truncate(time.Second).Format(timeFormat)
		res = append(res, []byte(tstr)...)
	}
	res = append(res, ',')
	if t.UpperBound.Time != nil {
		tstr := t.UpperBound.Time.UTC().Truncate(time.Second).Format(timeFormat)
		res = append(res, []byte(tstr)...)
	}
	if t.UpperBound.Inclusive {
		res = append(res, ']')
	} else {
		res = append(res, ')')
	}
	return res, nil
}

// Scan implements database/sql Scanner interface.
func (t *TSRange) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("TSRange.Scan: expected []byte, got %T (%q)", value, value)
	}
	if len(v) < 3 {
		return fmt.Errorf("TSRange.Scan: unexpected data %q", v)
	}
	if v[0] != '(' && v[0] != '[' {
		return fmt.Errorf("TSRange.Scan: unexpected data %q", v)
	}
	if v[len(v)-1] != ')' && v[len(v)-1] != ']' {
		return fmt.Errorf("TSRange.Scan: unexpected data %q", v)
	}
	if v[0] == '[' {
		t.LowerBound.Inclusive = true
	} else {
		t.LowerBound.Inclusive = false
	}
	commaIdx := bytes.IndexByte(v, ',')
	if commaIdx == -1 {
		return fmt.Errorf("TSRange.Scan: no comma %q", v)
	}
	lt := v[1:commaIdx]
	if len(lt) > 0 {
		lt = lt[1 : len(lt)-1]
		time, err := time.Parse(timeFormat, string(lt))
		if err != nil {
			return fmt.Errorf("TSRange.Scan: error parsing lower bound time %s: %s", lt, err)
		}
		t.LowerBound.Time = &time
	}
	ut := v[commaIdx+1 : len(v)-1]
	if len(ut) > 0 {
		ut = ut[1 : len(ut)-1]
		time, err := time.Parse(timeFormat, string(ut))
		if err != nil {
			return fmt.Errorf("TSRange.Scan: error parsing upper bound time %s: %s", ut, err)
		}
		t.UpperBound.Time = &time
	}
	if v[len(v)-1] == ']' {
		t.UpperBound.Inclusive = true
	} else {
		t.UpperBound.Inclusive = false
	}
	return nil
}

var (
	_ driver.Valuer = TSRange{}
	_ sql.Scanner   = &TSRange{}
)
