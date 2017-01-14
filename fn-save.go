package dbmdl

import (
	"errors"
	"reflect"
)

// Saves will add to the database or update an existing resource if a nonzero WHERE is provided
func Save(dlct, t string, target interface{}, where *WhereClause, fields ...string) error {

	// Set fields is not given already
	var ref = reflect.TypeOf(target)
	if len(fields) < 1 {
		for i := 0; i < ref.NumField(); i++ {
			field := ref.Field(i) // Get the field at index i
			if field.Tag.Get("dbmdl") == "" {
				continue
			}

			for _, tag := range getTagParameters(field) {
				if tag == omit {
					continue
				}
			}

			fields = append(fields, field.Name)
		}
	}

	// Build fieldsValues
	var fieldsValues = make(map[string]interface{})
	var val = reflect.ValueOf(target)
	for _, f := range fields {
		fieldsValues[f] = val.FieldByName(f).Interface()
	}

	// Handle dialects
	d, ok := dialects[dlct]
	if !ok {
		return errors.New("[dbmdl] Dialect " + dlct + " unknown")
	}

	// Handle query
	if len(where.Clauses) < 1 { // If the where clause is empty, INSERT:
		q := d.Insert(t, fieldsValues) // Build query
		query(nil, q...)               // Query, no return channel
	} else { // If the where clause is not empty, UPDATE:
		where.Dialect = d                     // Set the dialect of the WhereClause
		q := d.Update(t, fieldsValues, where) // Build query
		query(nil, q...)                      // Query, no return channel
	}

	return nil
}
