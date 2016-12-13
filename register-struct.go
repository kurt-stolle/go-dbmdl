package dbmdl

import (
	"errors"
	"log"
	"reflect"
)

// RegisterStruct registers a struct for use with dbmdl
func RegisterStruct(dlct string, t string, s interface{}) error {
	d, ok := dialects[dlct]
	if !ok {
		return errors.New("[dbmdl] Failed to register struct; dialect " + dlct + " unknown!")
	}

	refType := reflect.TypeOf(s).Elem()

	if _, exists := tables[refType]; exists {
		return errors.New("[dbmdl] Type " + refType.Name() + " is already registered!")
	}

	tables[refType] = &table{dialect: d, name: t}

	log.Println("[dbmdl] Registered struct: " + refType.Name())

	// Return possible errors from table creation
	return createTables(refType)
}
