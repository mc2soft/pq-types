package pq_types

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TypesSuite struct {
	db          *sql.DB
	skipPostGIS bool
}

var _ = Suite(&TypesSuite{})

func (s *TypesSuite) SetUpSuite(c *C) {
	db, err := sql.Open("postgres", "dbname=pq_types sslmode=disable")
	c.Assert(err, IsNil)
	s.db = db

	var version string
	row := db.QueryRow("SELECT version()")
	err = row.Scan(&version)
	c.Assert(err, IsNil)
	log.Print(version)

	s.db.Exec("DROP TABLE IF EXISTS pq_types")
	_, err = s.db.Exec(`CREATE TABLE pq_types(
		stringarray varchar[],
		int32_array int[],
		jsontext_varchar varchar,
		jsontext_json json,
		jsontext_jsonb jsonb
	)`)
	c.Check(err, IsNil)

	if _, err = s.db.Exec("SELECT PostGIS_full_version()"); err != nil {
		s.skipPostGIS = true
	} else {
		_, err = s.db.Exec("ALTER TABLE pq_types ADD COLUMN box box2d")
		c.Check(err, IsNil)

		_, err = s.db.Exec("SELECT AddGeometryColumn('pq_types','point','4326','POINT',2)")
		c.Check(err, IsNil)

		_, err = s.db.Exec("SELECT AddGeometryColumn('pq_types','polygon','4326','POLYGON',2)")
		c.Check(err, IsNil)
	}
}

func (s *TypesSuite) SetUpTest(c *C) {
	_, err := s.db.Exec("TRUNCATE TABLE pq_types")
	c.Check(err, IsNil)
}

func (s *TypesSuite) TearDownSuite(c *C) {
	s.db.Close()
}
