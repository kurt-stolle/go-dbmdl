package dbmdl

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
)

// Errors
var (
	ErrNotFound = sql.ErrNoRows
)

// Load will load a single struct from the database based on a where clause
// Target is a pointer to a struct
func Load(db *sql.DB, table string, target interface{}, where *WhereClause) error {
	// Check whether the dialect exists
	if where.Dialect == nil {
		return errors.New("WhereClause does not have a dialect set")
	}

	// First, verify whether the supplied target is actually a pointer
	var targetType = reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return errors.New("[dbmdl] target passed is not a pointer")
	}

	// Set references for later use
	targetType = targetType.Elem()
	targetValue := reflect.ValueOf(target).Elem()

	// Check whether we know of this type's existance
	if _, exists := tables[targetType]; !exists {
		return errors.New("[dbmdl] Type " + targetType.Name() + " is not a known type!")
	}

	// Get the fields
	var fields []string
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i) // Get the field at index i
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

	// Query using the same shit as Fetch Fields
	q, a := where.Dialect.FetchFields(table, fields, &Pagination{1, 1}, where)
	r, err := db.Query(q, a...)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		log.Panic(err)
	}

	// Create dummy variables that we can scan the results of the query into
	var addresses []interface{}
	for _, name := range fields {
		valField := targetValue.FieldByName(name)
		if !valField.CanAddr() {
			return errors.New("Field not found in array: " + name)
		}

		addresses = append(addresses, valField.Addr().Interface()) // Add the address of the field to the addresses array so that we can scan into this addresss later
	}

	// Wait for query to return a result and start scanning
	r.Scan(addresses...) // Scan into a by pointer targetTypeerence

	return nil
}
