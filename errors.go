package dbmdl

import (
	"database/sql"
	"errors"
)

// Errors
var (
	ErrNotFound       = sql.ErrNoRows
	ErrStructNotFound = errors.New("dbmdl: The struct provided is not registered with DBMDL")
	ErrNoDialect      = errors.New("dbmdl: WhereClause has no dialect set")
	ErrNoPointer      = errors.New("dbmdl: Target is not a pointer")
	ErrUnknownType    = errors.New("dbmdl: Unknown type requested")
)
