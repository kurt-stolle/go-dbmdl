package dbmdl

import (
	"database/sql"
	"reflect"
)

// Save will add to the database or update an existing resource if a nonzero WHERE is provided
func Save(db *sql.DB, t string, target interface{}, where *WhereClause, fields ...string) error {
	// Check dialect
	if where.Dialect == nil {
		return ErrNoDialect
	}

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

	// Create a channel for the reply
	res := make(chan *sql.Rows)
	defer close(res)

	// Handle query
	var q string
	var a []interface{}
	if len(where.Clauses) < 1 { // If the where clause is empty, INSERT:
		q, a = where.Dialect.Insert(t, fieldsValues) // Build query
	} else { // If the where clause is not empty, UPDATE:
		q, a = where.Dialect.Update(t, fieldsValues, where) // Build query
	}

	// Wait for response and close channel
	if _, err := db.Exec(q, a...); err != nil {
		return err
	}

	return nil
}
