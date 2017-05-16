package dbmdl

import (
	"database/sql"
	"log"
	"reflect"
)

// Load will load a single struct from the database based on a where clause
func (m *Model) Load(target interface{}, where *WhereClause) error {
	targetValue := reflect.ValueOf(target).Elem()

	// Get the fields
	fields := m.GetFields()

	// Query using the same shit as Fetch Fields
	q, a := m.Dialect.FetchFields(m.TableName, fields, NewPagination(1, 1), where)

	r := m.GetDatabase().QueryRow(q, a...)

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
