package dbmdl

import (
	"database/sql"
	"log"
	"reflect"
	"sync"
)

// Fetch loads data from a database, returns an array of interface, pagination is also updated automatically
func (m *Model) Fetch(pag *Pagination, where WhereSelector, fields ...string) ([]interface{}, error) {
	// If we did not supply and fields to be selected, select all fields
	if len(fields) < 1 {
		fields = m.GetFields()
	}

	// Do the following tasks concurrently
	var wg sync.WaitGroup
	var dummyVariables []reflect.Value // Slice to hold the values scanned from the *sql.Rows result
	var dummyVariablesAddresses []interface{}
	var data []interface{}

	// Build and execute the Query
	wg.Add(1)
	go func() {
		q, a := m.Dialect.FetchFields(m.TableName, fields, pag, where)

		r, err := m.GetDatabase().Query(q, a...)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
		defer r.Close()
		for _, name := range fields {
			f, found := m.Type.FieldByName(name)
			if !found {
				continue
			}

			var nwtyp = reflect.New(f.Type) // This returns a reflect value of a pointer to the type at f

			dummyVariables = append(dummyVariables, nwtyp.Elem())                        // Store our newly made Value
			dummyVariablesAddresses = append(dummyVariablesAddresses, nwtyp.Interface()) // Store the address as well, so that we can supply this to the Scan function later down the road
		}

		for r.Next() {
			var s = reflect.New(m.Type) // Create a new pointer to an empty struct of type targetType

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
		if err := pag.Load(m, where); err != nil {
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

	return data, nil
}

// FetchAny loads data from the database like Fetch, but without requiring a WhereClause
func (m *Model) FetchAny(pag *Pagination, fields ...string) ([]interface{}, error){
	return m.Fetch(pag, new(WhereClause), fields...)
}