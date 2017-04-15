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
		fields = getFields(targetType)
	}

	if len(fields) < 1 {
		return nil, nil, errors.New("Nothing to select, all fields flagged omit")
	}

	// Do the following tasks concurrently
	var wg sync.WaitGroup
	var dummyVariables []reflect.Value // Slice to hold the values scanned from the *sql.Rows result
	var dummyVariablesAddresses []interface{}
	var data []interface{}

	// Build and execute the Query
	wg.Add(1)
	go func() {
		q, a := d.FetchFields(t, fields, pag, where)

		r, err := db.Query(q, a...)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
		defer r.Close()
		for _, name := range fields {
			f, found := targetType.FieldByName(name)
			if !found {
				continue
			}

			var nwtyp = reflect.New(f.Type) // This returns a reflect value of a pointer to the type at f

			dummyVariables = append(dummyVariables, nwtyp.Elem())                        // Store our newly made Value
			dummyVariablesAddresses = append(dummyVariablesAddresses, nwtyp.Interface()) // Store the address as well, so that we can supply this to the Scan function later down the road
		}

		for r.Next() {
			var s = reflect.New(targetType) // Create a new pointer to an empty struct of type targetType

			r.Scan(dummyVariablesAddresses...) // Scan into the slice we populated with dummy variables earlier

			for i, v := range dummyVariables {
				s.Elem().FieldByName(fields[i]).Set(v) // Set values in our new struct
			}

			data = append(data, s.Interface()) // Append the interface value of the pointer to the previously created targetType type.
		}

		wg.Done()
	}()

	// Pagination
	wg.Add(1)
	go func() {
		if err := pag.Load(db, t, where); err != nil {
			if err == sql.ErrNoRows {
				pag.First = 1
				pag.Next = 1
				pag.Prev = 1
				pag.Last = 1
				return
			}
			panic(err)
		}
		wg.Done()
	}()

	// Wait
	wg.Wait()

	return data, pag, nil
}
