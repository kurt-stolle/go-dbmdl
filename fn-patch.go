package dbmdl

import (
	"database/sql"
	"reflect"
)

// Patch works much the same as Save, but only performs an update according to a map.
// Patch is used when there is no struct initialization required, e.g. when only fields in the database need to be updated.
func Patch(db *sql.DB, t string, sRef interface{}, where *WhereClause, update map[string]interface{}) error {
	// First, verify whether the supplied target is actually a pointer
	var targetType = reflect.TypeOf(sRef)
	if targetType.Kind() != reflect.Ptr {
		return ErrNoPointer
	}
	targetType = targetType.Elem()

	// Create a new struct so that we can pass this to Save
	var newStruct = reflect.New(targetType)
	var fields []string

	// Verify the map's entries
	for k, v := range update {
		if _, found := targetType.FieldByName(k); !found {
			return ErrNotFound
		}

		field := newStruct.FieldByName(k)
		field.Set(reflect.ValueOf(v))

		fields = append(fields, k)
	}

	// Perform Save
	return Save(db, t, newStruct.Addr().Interface(), where, fields...)
}
