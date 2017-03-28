package dbmdl

import (
	"database/sql"
	"errors"
	"reflect"
)

// getDT is a helper function for fetching the dialect and table
func getDT(typ reflect.Type) (*Dialect, string, error) {
	if dt, exists := tables[typ]; exists {
		return dt.dialect, dt.name, nil
	}

	return nil, "", ErrStructNotFound
}

// RegisterStruct registers a struct for use with dbmdl
func RegisterStruct(db *sql.DB, dlct string, t string, s interface{}) error {
	d, ok := dialects[dlct]
	if !ok {
		return ErrNoDialect
	}

	// s Must be a pointer to a struct of the type we are looking for
	// It is advised but not required to use the nil value
	// because this saves us memory
	refType := getReflectType(s)

	// A struct may only be registered once. Each project should define its own structs
	if _, exists := tables[refType]; exists {
		return errors.New("dbmdl: Struct is already registered")
	}

	tables[refType] = &table{dialect: d, name: t}

	// Return possible errors from table creation
	return nil
}
