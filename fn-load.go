package dbmdl

import (
	"database/sql"
	"log"
	"reflect"
)

// Load will load a single struct from the database based on a where clause
// Target is a pointer to a struct
func Load(db *sql.DB, target interface{}, where *WhereClause) error {
	// Check whether the dialect exists
	if where.Dialect == nil {
		return ErrNoDialect
	}

	// First, verify whether the supplied target is actually a pointer
	var targetType = reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		panic(ErrNoPointer)
	}

	// Set references for later use
	targetType = targetType.Elem()
	targetValue := reflect.ValueOf(target).Elem()

	// Get the dialect and table name
	d, t, err := getDT(targetType)
	if err != nil {
		return err
	}

	// Check whether we know of this type's existance
	if _, exists := tables[targetType]; !exists {
		return ErrUnknownType
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
	q, a := d.FetchFields(t, fields, NewPagination(1, 1), where)

	r := db.QueryRow(q, a...)

	// Create dummy variables that we can scan the results of the query into
	var addresses []interface{}
	for _, name := range fields {
		valField := targetValue.FieldByName(name)
		if !valField.CanAddr() {
			log.Panic("dbmdl: Field '" + name + "' not found")
		}

		addresses = append(addresses, valField.Addr().Interface()) // Add the address of the field to the addresses array so that we can scan into this addresss later
	}

	// Wait for query to return a result and start scanning
	if err := r.Scan(addresses...); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		log.Panic("dbmdl: ", err)
	}

	return nil
}
