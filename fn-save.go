package dbmdl

import (
	"database/sql"
	"errors"
	"reflect"
)

// Save will add to the database or update an existing resource if a nonzero WHERE is provided
func Save(dlct, t string, target interface{}, where *WhereClause, fields ...string) error {

	// Set fields is not given already
	var ref = reflect.TypeOf(target)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

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
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, f := range fields {
		fieldsValues[f] = val.FieldByName(f).Interface()
	}

	// Handle dialects
	d, ok := dialects[dlct]
	if !ok {
		return errors.New("[dbmdl] Dialect " + dlct + " unknown")
	}

	// Create a channel for the reply
	res := make(chan *sql.Rows)
	defer close(res)

	// Handle query
	if len(where.Clauses) < 1 { // If the where clause is empty, INSERT:
		q := d.Insert(t, fieldsValues) // Build query
		query(res, q...)               // Query, no return channel
	} else { // If the where clause is not empty, UPDATE:
		where.Dialect = d                     // Set the dialect of the WhereClause
		q := d.Update(t, fieldsValues, where) // Build query
		query(res, q...)                      // Query, no return channel
	}

	// Wait for response and close channel
	r := <-res
	defer r.Close()

	return nil
}
