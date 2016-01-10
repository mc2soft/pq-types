# pq-types [![Build Status](https://travis-ci.org/mc2soft/pq-types.svg?branch=master)](https://travis-ci.org/mc2soft/pq-types) [![GoDoc](https://godoc.org/github.com/mc2soft/pq-types?status.svg)](http://godoc.org/github.com/mc2soft/pq-types)

This Go package provides additional types for PostgreSQL:

* `Int32Array` for `int[]` (compatible with [`intarray`](http://www.postgresql.org/docs/current/static/intarray.html) module);
* `StringArray` for `varchar[]`.
* `JSONText` for `varchar`, `text`, `json` and `jsonb`.

Install it: `go get github.com/mc2soft/pq-types`
