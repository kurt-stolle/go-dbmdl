package dbmdl

import (
	"reflect"
)

// Patch works much the same as Save, but only performs an update according to a map.
// Patch is used when there is no struct initialization required, e.g. when only fields in the database need to be updated.
func (m *Model) Patch(where WhereSelector, update map[string]interface{}) error {
	// Create a new struct so that we can pass this to Save
	var typePtr = reflect.New(m.Type)
	var fields []string

	// Verify the map's entries
	for k, v := range update {
		if _, found := m.Type.FieldByName(k); !found {
			continue // Skip fields that don't exist
		}

		field := typePtr.Elem().FieldByName(k)
		field.Set(reflect.ValueOf(v))

		fields = append(fields, k)
	}

	// Check whether we've actually updated anything
	if len(fields) < 1 {
		return nil
	}

	// Perform Save
	return m.Save(typePtr.Interface(), where, fields...)
}
