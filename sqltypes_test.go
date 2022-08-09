package pq_types

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	. "gopkg.in/check.v1"
)

type Logger interface {
	Logf(format string, args ...interface{})
}

type DB struct {
	*sql.DB
	l Logger
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db.l != nil {
		db.l.Logf("%s (args = %#v)", query, args)
	}
	return db.DB.Query(query, args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	if db.l != nil {
		db.l.Logf("%s (args = %#v)", query, args)
	}
	return db.DB.QueryRow(query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db.l != nil {
		db.l.Logf("%s (args = %#v)", query, args)
	}
	return db.DB.Exec(query, args...)
}

func Test(t *testing.T) { TestingT(t) }

type TypesSuite struct {
	db          *DB
	skipJSON    bool
	skipJSONB   bool
	skipPostGIS bool
	skipTSRange bool
}

var _ = Suite(&TypesSuite{})

func (s *TypesSuite) SetUpSuite(c *C) {
	db, err := sql.Open("postgres", "dbname=pq_types sslmode=disable")
	c.Assert(err, IsNil)
	s.db = &DB{
		DB: db,
		l:  c,
	}

	// log full version
	var version string
	row := db.QueryRow("SELECT version()")
	err = row.Scan(&version)
	c.Assert(err, IsNil)
	log.Print(version)

	// check minor version
	row = db.QueryRow("SHOW server_version")
	err = row.Scan(&version)
	c.Assert(err, IsNil)
	minor, err := strconv.Atoi(strings.Split(version, ".")[1])
	c.Assert(err, IsNil)

	// check json and jsonb support
	if minor <= 1 {
		log.Print("json not available")
		s.skipJSON = true
		log.Print("tsrange not available")
		s.skipTSRange = true
	}
	if minor <= 3 {
		log.Print("jsonb not available")
		s.skipJSONB = true
	}

	s.db.Exec("DROP TABLE IF EXISTS pq_types")
	_, err = s.db.Exec(`CREATE TABLE pq_types(
		string_array varchar[],
		int32_array int[],
		int64_array bigint[],
		jsontext_varchar varchar,
		null_str varchar,
		null_int32 int4,
		null_int64 int8,
		null_timestamp timestamptz
	)`)
	c.Assert(err, IsNil)

	if !s.skipJSON {
		_, err = s.db.Exec(`ALTER TABLE pq_types ADD COLUMN jsontext_json json`)
		c.Assert(err, IsNil)
	}

	if !s.skipJSONB {
		_, err = s.db.Exec(`ALTER TABLE pq_types ADD COLUMN jsontext_jsonb jsonb`)
		c.Assert(err, IsNil)
	}

	if !s.skipTSRange {
		_, err = s.db.Exec(`ALTER TABLE pq_types ADD COLUMN tsrange tsrange`)
		c.Assert(err, IsNil)
	}

	// check PostGIS
	db.Exec("CREATE EXTENSION postgis")
	row = db.QueryRow("SELECT PostGIS_full_version()")
	err = row.Scan(&version)
	if err == nil {
		log.Print(version)

		_, err = s.db.Exec(`ALTER TABLE pq_types
			ADD COLUMN point geography(POINT, 4326),
			ADD COLUMN box box2d,
			ADD COLUMN polygon geography(POLYGON, 4326)
		`)
		c.Assert(err, IsNil)
	} else {
		log.Printf("PostGIS not available: %s", err)
		s.skipPostGIS = true
	}
}

func (s *TypesSuite) SetUpTest(c *C) {
	s.db.l = c
	_, err := s.db.Exec("TRUNCATE TABLE pq_types")
	c.Check(err, IsNil)
}

func (s *TypesSuite) TearDownSuite(c *C) {
	s.db.l = c
	s.db.Close()
}

func (s *TypesSuite) TestEmpty(c *C) {
	type record struct {
		i32a Int32Array
		i64a Int64Array
		sa   StringArray
	}

	for _, r := range []record{
		{},
		{i32a: Int32Array{}, i64a: Int64Array{}, sa: StringArray{}},
	} {
		s.SetUpTest(c)

		_, err := s.db.Exec(
			"INSERT INTO pq_types (int32_array, int64_array, string_array) VALUES($1, $2, $3)",
			r.i32a, r.i64a, r.sa,
		)
		c.Assert(err, IsNil)

		var r1 record
		row := s.db.QueryRow("SELECT int32_array, int64_array, string_array FROM pq_types")
		err = row.Scan(&r1.i32a, &r1.i64a, &r1.sa)
		c.Check(err, IsNil)
		c.Check(r1, DeepEquals, r)
	}
}
