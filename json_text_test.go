package pq_types

import (
	"fmt"
	. "gopkg.in/check.v1"
)

func (s *TypesSuite) TestJSONText(c *C) {
	type testData struct {
		j JSONText
		b []byte
	}
	for _, d := range []testData{
		{JSONText(nil), []byte(nil)},
		{JSONText(`null`), []byte(`null`)},
		{JSONText(`{}`), []byte(`{}`)},
		{JSONText(`[]`), []byte(`[]`)},
		{JSONText(`[{"b": true, "n": 123}, {"s": "foo", "obj": {"f1": 456, "f2": false}}, [null]]`),
			[]byte(`[{"b": true, "n": 123}, {"s": "foo", "obj": {"f1": 456, "f2": false}}, [null]]`)},
	} {
		s.SetUpTest(c)

		_, err := s.db.Exec("INSERT INTO pq_types (jsontext_varchar, jsontext_json, jsontext_jsonb) VALUES($1, $2, $3)",
			d.j, d.j, d.j)
		c.Assert(err, IsNil)

		for _, col := range []string{"jsontext_varchar", "jsontext_json", "jsontext_jsonb"} {
			b1 := []byte(`{"foo": "bar"}`)
			j1 := JSONText(`{"foo": "bar"}`)
			err = s.db.QueryRow(fmt.Sprintf("SELECT %s, %s FROM pq_types", col, col)).Scan(&b1, &j1)
			c.Check(err, IsNil)
			c.Check(b1, DeepEquals, d.b, Commentf("\nb1  = %#q\nd.b = %#q", b1, d.b))
			c.Check(j1, DeepEquals, d.j)
		}
	}
}
