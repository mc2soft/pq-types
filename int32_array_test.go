package pq_types

import (
	"database/sql"
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *TypesSuite) TestInt32Array(c *C) {
	type testData struct {
		a Int32Array
		b []byte
	}
	for _, d := range []testData{
		{Int32Array(nil), []byte(nil)},
		{Int32Array{}, []byte(`{}`)},
		{Int32Array{1}, []byte(`{1}`)},
		{Int32Array{1, 0, -3}, []byte(`{1,0,-3}`)},
		{Int32Array{-3, 0, 1}, []byte(`{-3,0,1}`)},
	} {
		s.SetUpTest(c)

		_, err := s.db.Exec("INSERT INTO pq_types (int32_array) VALUES($1)", d.a)
		c.Assert(err, IsNil)

		b1 := []byte("42")
		a1 := Int32Array{42}
		err = s.db.QueryRow("SELECT int32_array, int32_array FROM pq_types").Scan(&b1, &a1)
		c.Check(err, IsNil)
		c.Check(b1, DeepEquals, d.b, Commentf("\nb1  = %#q\nd.b = %#q", b1, d.b))
		c.Check(a1, DeepEquals, d.a)

		// check db array length
		var length sql.NullInt64
		err = s.db.QueryRow("SELECT array_length(int32_array, 1) FROM pq_types").Scan(&length)
		c.Check(err, IsNil)
		c.Check(length.Valid, Equals, len(d.a) > 0)
		c.Check(length.Int64, Equals, int64(len(d.a)))

		// check db array elements
		for i := 0; i < len(d.a); i++ {
			q := fmt.Sprintf("SELECT int32_array[%d] FROM pq_types", i+1)
			var el sql.NullInt64
			err = s.db.QueryRow(q).Scan(&el)
			c.Check(err, IsNil)
			c.Check(el.Valid, Equals, true)
			c.Check(el.Int64, Equals, int64(d.a[i]))
		}
	}
}

func (s *TypesSuite) TestInt32ArrayEqualWithoutOrder(c *C) {
	c.Check(Int32Array{1, 0, -3}.EqualWithoutOrder(Int32Array{-3, 0, 1}), Equals, true)
	c.Check(Int32Array{1, 0, -3}.EqualWithoutOrder(Int32Array{1}), Equals, false)
	c.Check(Int32Array{1, 0, -3}.EqualWithoutOrder(Int32Array{1, 0, 42}), Equals, false)
	c.Check(Int32Array{}.EqualWithoutOrder(Int32Array{}), Equals, true)
	c.Check(Int32Array{}.EqualWithoutOrder(Int32Array{1}), Equals, false)
}
