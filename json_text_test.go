package pq_types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
		b1, err := json.Marshal(d.j)
		c.Check(err, IsNil)
		b := bytes.Replace(d.b, []byte(` `), nil, -1)
		if d.j == nil {
			// special case
			c.Check(b1, DeepEquals, []byte(`null`))
		} else {
			c.Check(b1, DeepEquals, b, Commentf("\nb1  = %#q\nb  = %#q", b1, b))
		}

		for _, col := range []string{"jsontext_varchar", "jsontext_json", "jsontext_jsonb"} {
			if strings.HasSuffix(col, "json") && s.skipJSON {
				continue
			}
			if strings.HasSuffix(col, "jsonb") && s.skipJSONB {
				continue
			}

			s.SetUpTest(c)

			_, err = s.db.Exec(fmt.Sprintf("INSERT INTO pq_types (%s) VALUES($1)", col), d.j)
			c.Assert(err, IsNil)

			b1 := []byte(`{"foo": "bar"}`)
			j1 := JSONText(`{"foo": "bar"}`)
			err = s.db.QueryRow(fmt.Sprintf("SELECT %s, %s FROM pq_types", col, col)).Scan(&b1, &j1)
			c.Check(err, IsNil)
			c.Check(b1, DeepEquals, d.b, Commentf("\nb1  = %#q\nd.b = %#q", b1, d.b))
			c.Check(j1, DeepEquals, d.j)
		}
	}

	for _, j := range []JSONText{
		JSONText{},
	} {
		for _, col := range []string{"jsontext_varchar", "jsontext_json", "jsontext_jsonb"} {
			if strings.HasSuffix(col, "json") && s.skipJSON {
				continue
			}
			if strings.HasSuffix(col, "jsonb") && s.skipJSONB {
				continue
			}

			s.SetUpTest(c)

			_, err := s.db.Exec(fmt.Sprintf("INSERT INTO pq_types (%s) VALUES($1)", col), j)
			c.Check(err, DeepEquals, errors.New(`sql: converting Exec argument #0's type: unexpected end of JSON input`))
		}
	}
}
