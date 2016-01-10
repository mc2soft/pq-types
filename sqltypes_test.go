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

func Test(t *testing.T) { TestingT(t) }

type TypesSuite struct {
	db          *sql.DB
	skipJSON    bool
	skipJSONB   bool
	skipPostGIS bool
}

var _ = Suite(&TypesSuite{})

func (s *TypesSuite) SetUpSuite(c *C) {
	db, err := sql.Open("postgres", "dbname=pq_types sslmode=disable")
	c.Assert(err, IsNil)
	s.db = db

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
	}
	if minor <= 3 {
		log.Print("jsonb not available")
		s.skipJSONB = true
	}

	s.db.Exec("DROP TABLE IF EXISTS pq_types")
	_, err = s.db.Exec(`CREATE TABLE pq_types(
		stringarray varchar[],
		int32_array int[],
		jsontext_varchar varchar
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

	// check PostGIS
	_, err = db.Exec("CREATE EXTENSION postgis")
	if err == nil {
		row = db.QueryRow("SELECT PostGIS_full_version()")
		err = row.Scan(&version)
	}
	if err == nil {
		log.Print(version)

		_, err = s.db.Exec("ALTER TABLE pq_types ADD COLUMN box box2d")
		c.Check(err, IsNil)

		_, err = s.db.Exec("SELECT AddGeometryColumn('pq_types','point','4326','POINT',2)")
		c.Check(err, IsNil)

		_, err = s.db.Exec("SELECT AddGeometryColumn('pq_types','polygon','4326','POLYGON',2)")
		c.Check(err, IsNil)
	} else {
		log.Printf("PostGIS not available: %s", err)
		s.skipPostGIS = true
	}
}

func (s *TypesSuite) SetUpTest(c *C) {
	_, err := s.db.Exec("TRUNCATE TABLE pq_types")
	c.Check(err, IsNil)
}

func (s *TypesSuite) TearDownSuite(c *C) {
	s.db.Close()
}
