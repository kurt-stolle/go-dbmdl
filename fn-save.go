package dbmdl

import (
	"database/sql"
	"reflect"
)

// Save will add to the database or update an existing resource if a nonzero WHERE is provided
func Save(db *sql.DB, target interface{}, where *WhereClause, fields ...string) error {
	// Set fields is not given already
	var targetType = reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		panic(ErrNoPointer)
	}
	targetType = targetType.Elem()

	// If there are no fields provided, select every field without an omit tag
	if len(fields) < 1 {
		for i := 0; i < targetType.NumField(); i++ {
			field := targetType.Field(i) // Get the field at index i
			if field.Tag.Get("dbmdl") == "" {
				continue
			}

			for _, tag := range getTagParameters(field) {
				if tag == omit || regExtern.MatchString(tag) {
					continue
				}
			}

			fields = append(fields, field.Name)
		}
	}

	// Build fieldsValues
	var fieldsValues = make(map[string]interface{})
	var val = reflect.ValueOf(target).Elem()

	// Get the dialect and table name
	d, t, err := getDT(reflect.TypeOf(target).Elem())
	if err != nil {
		return err
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
	if where == nil || len(where.Clauses) < 1 { // If the where clause is empty, INSERT:
		q, a = d.Insert(t, fieldsValues) // Build query
	} else { // If the where clause is not empty, UPDATE:
		q, a = d.Update(t, fieldsValues, where) // Build query
	}

	// Wait for response and close channel
	if _, err := db.Exec(q, a...); err != nil {
		return err
	}

	return nil
}
