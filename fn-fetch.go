package dbmdl

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"sync"
)

// Fetch loads data from a database a populates the struct
// sRef is a pointer to the struct, only used for getting the reflection type
func Fetch(db *sql.DB, sRef interface{}, where *WhereClause, pag *Pagination, fields ...string) ([]interface{}, *Pagination, error) {

	// Set the reference, but check whether it's a pointer first
	targetType := getReflectType(sRef)

	// Get the dialect and table name
	d, t, err := getDT(targetType)
	if err != nil {
		return nil, nil, err
	}

	// Check whether we know of this type's existance
	if _, exists := tables[targetType]; !exists {
		return nil, nil, ErrUnknownType
	}

	// Fallbacks
	if pag == nil {
		pag = NewPagination(1, 1)
	}

	if where == nil {
		where = NewWhereClause(d)
	}

	// If we did not supply and fields to be selected, select all fields
	if len(fields) < 1 {
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
	}

	if len(fields) < 1 {
		return nil, nil, errors.New("Nothing to select, all fields flagged omit")
	}

	// Do the following tasks concurrently
	var wg sync.WaitGroup
	var r *sql.Rows
	var dummyVariables []reflect.Value // Slice to hold the values scanned from the *sql.Rows result
	var dummyVariablesAddresses []interface{}
	var data []interface{}

	// Build and execute the Query
	wg.Add(1)
	go func() {
		q, a := d.FetchFields(t, fields, pag, where)

		rows, err := db.Query(q, a...)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}

		r = rows // Make rows available outside this scope

		wg.Done()
	}()

	// Results
	wg.Add(1)
	go func() {
		// Create dummy variables that we can scan the results of the query into
		for _, name := range fields {
			f, found := targetType.FieldByName(name)
			if !found {
				continue
			}

			var nwtyp = reflect.New(f.Type) // This returns a reflect value of a pointer to the type at f

			dummyVariables = append(dummyVariables, nwtyp.Elem())                        // Store our newly made Value
			dummyVariablesAddresses = append(dummyVariablesAddresses, nwtyp.Interface()) // Store the address as well, so that we can supply this to the Scan function later down the road
		}
		wg.Done()
	}()

	// Pagination
	wg.Add(1)
	go func() {
		pag.Load(db, t, where)
		wg.Done()
	}()

	// Wait
	wg.Wait()

	defer r.Close()

	// Iterate over rows found
	for r.Next() {
		var s = reflect.New(targetType) // Create a new pointer to an empty struct of type targetType

		r.Scan(dummyVariablesAddresses...) // Scan into the slice we populated with dummy variables earlier

		for i, v := range dummyVariables {
			s.Elem().FieldByName(fields[i]).Set(v) // Set values in our new struct
		}

		data = append(data, s.Interface()) // Append the interface value of the pointer to the previously created targetType type.
	}

	return data, pag, nil
}
