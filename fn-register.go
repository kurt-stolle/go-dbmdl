package dbmdl

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
)

// RegisterStruct registers a struct for use with dbmdl
func RegisterStruct(db *sql.DB, dlct string, t string, s interface{}) error {
	d, ok := dialects[dlct]
	if !ok {
		return ErrNoDialect
	}

	refType := reflect.TypeOf(s).Elem()

	if _, exists := tables[refType]; exists {
		return errors.New("dbmdl: Type " + refType.Name() + " is already registered!")
	}

	tables[refType] = &table{dialect: d, name: t}

	log.Println("Registered struct: " + refType.Name())

	// Return possible errors from table creation
	return createTables(db, refType)
}
