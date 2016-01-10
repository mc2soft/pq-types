package pq_types

import (
	"database/sql"
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *TypesSuite) TestStringArray(c *C) {
	type testData struct {
		a StringArray
		b []byte
	}
	for _, d := range []testData{
		{StringArray(nil), []byte(nil)},
		{StringArray{}, []byte(`{}`)},

		{StringArray{`1234567`}, []byte(`{1234567}`)},
		{StringArray{`abc123, def456 xyz789`, `абв`, `世界,`}, []byte(`{"abc123, def456 xyz789",абв,"世界,"}`)},

		{StringArray{"", "`", "``", "```", "````"}, []byte("{\"\",`,``,```,````}")},
		{StringArray{``, `'`, `''`, `'''`, `''''`}, []byte(`{"",','',''',''''}`)},
		{StringArray{``, `"`, `""`, `"""`, `""""`}, []byte(`{"","\"","\"\"","\"\"\"","\"\"\"\""}`)},
		{StringArray{``, `,`, `,,`, `,,,`, `,,,,`}, []byte(`{"",",",",,",",,,",",,,,"}`)},
		{StringArray{``, `\`, `\\`, `\\\`, `\\\\`}, []byte(`{"","\\","\\\\","\\\\\\","\\\\\\\\"}`)},
		{StringArray{``, `{`, `{{`, `}}`, `}`, `{{}}`}, []byte(`{"","{","{{","}}","}","{{}}"}`)},

		{StringArray{`\{`, `\\{{`, `\}\}`, `\}}`}, []byte(`{"\\{","\\\\{{","\\}\\}","\\}}"}`)},
		{StringArray{`\"'`, `\\"`, `\\\"`, `"\"\\""`}, []byte(`{"\\\"'","\\\\\"","\\\\\\\"","\"\\\"\\\\\"\""}`)},
	} {
		s.SetUpTest(c)

		_, err := s.db.Exec("INSERT INTO pq_types (stringarray) VALUES($1)", d.a)
		c.Assert(err, IsNil)

		b1 := []byte("lalala")
		a1 := StringArray{"lalala"}
		err = s.db.QueryRow("SELECT stringarray, stringarray FROM pq_types").Scan(&b1, &a1)
		c.Check(err, IsNil)
		c.Check(b1, DeepEquals, d.b, Commentf("\nb1  = %#q\nd.b = %#q", b1, d.b))
		c.Check(a1, DeepEquals, d.a)

		// check db array length
		var length sql.NullInt64
		err = s.db.QueryRow("SELECT array_length(stringarray, 1) FROM pq_types").Scan(&length)
		c.Check(err, IsNil)
		c.Check(length.Valid, Equals, len(d.a) > 0)
		c.Check(length.Int64, Equals, int64(len(d.a)))

		// check db array elements
		for i := 0; i < len(d.a); i++ {
			q := fmt.Sprintf("SELECT stringarray[%d] FROM pq_types", i+1)
			var el sql.NullString
			err = s.db.QueryRow(q).Scan(&el)
			c.Check(err, IsNil)
			c.Check(el.Valid, Equals, true)
			c.Check(el.String, Equals, d.a[i])
		}
	}
}
