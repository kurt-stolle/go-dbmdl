package dbmdl

import (
	"database/sql"
	"errors"
)

// Errors
var (
	ErrNotFound    = sql.ErrNoRows
	ErrNoDialect   = errors.New("dbmdl: WhereClause has no dialect set")
	ErrNoPointer   = errors.New("dbmdl: Target is not a pointer")
	ErrUnknownType = errors.New("dbmdl: Unknown type requested")
)
