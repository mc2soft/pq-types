package pq_types

import (
	"database/sql"
	"time"
)

// NullString covers trivial case of string to sql.NullString conversion assuming empty string to be NULL
func NullString(src string) sql.NullString {
	return sql.NullString{
		String: src,
		Valid:  src != "",
	}
}

// NullInt32 covers trivial case of int32 to sql.NullInt32 conversion assuming 0 to be NULL
func NullInt32(src int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: src,
		Valid: src != 0,
	}
}

// NullInt64 covers trivial case of int64 to sql.NullInt64 conversion assuming 0 to be NULL
func NullInt64(src int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: src,
		Valid: src != 0,
	}
}

// NullTimestampP converts *time.Time to a sql.NullTime
func NullTimestampP(src *time.Time) sql.NullTime {
	if src == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  *src,
		Valid: true,
	}
}
