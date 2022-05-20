package pq_types

import (
	"database/sql"
	"fmt"
	"time"

	"gopkg.in/check.v1"
)

func insertQuery(col string) string {
	return fmt.Sprintf("INSERT INTO pq_types (%s) VALUES($1)", col)
}

func selectQuery(col string) string {
	return fmt.Sprintf("SELECT %s FROM pq_types LIMIT 1;", col)
}

// For each type there need to be a pair of tests: with empty -> nil value and with non-empty value -> value.

func (s *TypesSuite) TestConversionNullString(c *check.C) {
	cases := []sql.NullString{
		{String: ""},
		{String: "truly random string", Valid: true},
	}
	for _, expected := range cases {
		val := NullString(expected.String)
		_, err := s.db.Exec(insertQuery("null_str"), val)
		c.Assert(err, check.IsNil)

		var actual sql.NullString
		err = s.db.QueryRow(selectQuery("null_str")).Scan(actual)
		c.Check(err, check.IsNil)
		c.Check(actual, check.DeepEquals, expected)
	}
}

func (s *TypesSuite) TestConversionNullInt32(c *check.C) {
	cases := []sql.NullInt32{
		{Int32: 0},
		{Int32: 0xabc, Valid: true},
	}
	for _, expected := range cases {
		val := NullInt32(expected.Int32)
		_, err := s.db.Exec(insertQuery("null_int32"), val)
		c.Assert(err, check.IsNil)

		var actual sql.NullInt32
		err = s.db.QueryRow(selectQuery("null_int32")).Scan(actual)
		c.Check(err, check.IsNil)
		c.Check(actual, check.DeepEquals, expected)
	}
}

func (s *TypesSuite) TestConversionNullInt64(c *check.C) {
	cases := []sql.NullInt64{
		{Int64: 0},
		{Int64: 0xabcdef, Valid: true},
	}
	for _, expected := range cases {
		val := NullInt64(expected.Int64)
		_, err := s.db.Exec(insertQuery("null_int64"), val)
		c.Assert(err, check.IsNil)

		var actual sql.NullInt64
		err = s.db.QueryRow(selectQuery("null_int64")).Scan(actual)
		c.Check(err, check.IsNil)
		c.Check(actual, check.DeepEquals, expected)
	}
}

func (s *TypesSuite) TestConversionNullTimestamp(c *check.C) {
	// here we use another approach, as input is a pointer
	now := time.Now()
	cases := []*time.Time{nil, &now}

	for _, expected := range cases {
		s.SetUpTest(c)

		val := NullTimestampP(expected)

		_, err := s.db.Exec(insertQuery("null_timestamp"), val)
		c.Assert(err, check.IsNil)

		var actual *time.Time
		err = s.db.QueryRow(selectQuery("null_timestamp")).Scan(&actual)
		c.Check(err, check.IsNil)
		c.Check(actual, check.Equals, expected)
	}
}
