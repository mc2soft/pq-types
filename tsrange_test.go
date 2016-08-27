package pq_types

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"
)

func (s *TypesSuite) TestTSRange(c *C) {
	type testData struct {
		ts TSRange
		s  string
	}
	upperTime := time.Now().UTC().Truncate(time.Second)
	lowerTime := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	utStr := upperTime.Format(timeFormat)
	ltStr := lowerTime.Format(timeFormat)
	for _, d := range []testData{
		{TSRange{TimeBound{true, &lowerTime}, TimeBound{true, &upperTime}}, fmt.Sprintf(`["%s","%s"]`, ltStr, utStr)},
		{TSRange{TimeBound{false, &lowerTime}, TimeBound{false, &upperTime}}, fmt.Sprintf(`("%s","%s")`, ltStr, utStr)},
		{TSRange{TimeBound{false, &lowerTime}, TimeBound{true, &upperTime}}, fmt.Sprintf(`("%s","%s"]`, ltStr, utStr)},
		{TSRange{TimeBound{true, &lowerTime}, TimeBound{false, &upperTime}}, fmt.Sprintf(`["%s","%s")`, ltStr, utStr)},
		{TSRange{TimeBound{false, nil}, TimeBound{true, &upperTime}}, fmt.Sprintf(`(,"%s"]`, utStr)},
		{TSRange{TimeBound{true, &lowerTime}, TimeBound{false, nil}}, fmt.Sprintf(`["%s",)`, ltStr)},
		{TSRange{TimeBound{false, nil}, TimeBound{false, &upperTime}}, fmt.Sprintf(`(,"%s")`, utStr)},
		{TSRange{TimeBound{false, &lowerTime}, TimeBound{false, nil}}, fmt.Sprintf(`("%s",)`, ltStr)},
		{TSRange{TimeBound{false, nil}, TimeBound{false, nil}}, "(,)"},
	} {
		s.SetUpTest(c)
		_, err := s.db.Exec("INSERT INTO pq_types (tsrange) VALUES($1)", d.ts)
		c.Assert(err, IsNil)

		var el TSRange
		var els string
		err = s.db.QueryRow("SELECT tsrange, tsrange FROM pq_types").Scan(&el, &els)
		c.Check(err, IsNil)
		c.Check(d.ts, DeepEquals, el)
		c.Check(d.s, Equals, els)
	}
}
