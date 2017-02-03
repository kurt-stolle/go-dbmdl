package dbmdl

import (
	"database/sql"
	"errors"
	"reflect"
)

// Result holds an array of structs and a pagination object
type Result struct {
	Data       []interface{}
	Pagination *Pagination
}

// Fetch loads data from a database a populates the struct
// sRef is a pointer to the struct, only used for getting the reflection type
func Fetch(t string, sRef interface{}, where *WhereClause, pag *Pagination, fields ...string) (*Result, error) {
	// Check whether the dialect exists
	if where.Dialect == nil {
		return nil, errors.New("Invalid dialect")
	}

	// Set the targetTypeerence, but check whether it's a pointer first
	targetType := reflect.TypeOf(sRef)
	if targetType.Kind() != reflect.Ptr {
		return nil, errors.New("[dbmdl] target passed is not a pointer")
	}
	targetType = targetType.Elem()

	// Check whether we know of this type's existance
	if _, exists := tables[targetType]; !exists {
		return nil, errors.New("[dbmdl] Type " + targetType.Name() + " is not a known type!")
	}

	// Fallbacks
	if pag == nil {
		pag = &Pagination{1, 1}
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
		return nil, errors.New("We have nothing to select because all fields are ommited by dbmdl")
	}

	// Build and execute the Query
	q := where.Dialect.FetchFields(t, fields, pag, where)

	c := make(chan *sql.Rows)
	query(c, q...)

	// Results
	var res = &Result{}
	res.Pagination = pag

	// Create dummy variables that we can scan the results of the query into
	var dummyVariables []reflect.Value // Slice to hold the values scanned from the *sql.Rows result
	var dummyVariablesAddresses []interface{}
	for _, name := range fields {
		f, found := targetType.FieldByName(name)
		if !found {
			return nil, errors.New("Field not found in array: " + name)
		}

		var nwtyp = reflect.New(f.Type) // This returns a reflect value of a pointer to the type at f

		dummyVariables = append(dummyVariables, nwtyp.Elem())                        // Store our newly made Value
		dummyVariablesAddresses = append(dummyVariablesAddresses, nwtyp.Interface()) // Store the address as well, so that we can supply this to the Scan function later down the road
	}

	// Wait for the channel to return rows
	r := <-c
	defer r.Close()
	defer close(c)

	// Iterate over rows found
	for r.Next() {
		var s = reflect.New(targetType) // Create a new pointer to an empty struct of type targetType

		r.Scan(dummyVariablesAddresses...) // Scan into the slice we populated with dummy variables earlier

		for i, v := range dummyVariables {
			s.Elem().FieldByName(fields[i]).Set(v) // Set values in our new struct
		}

		res.Data = append(res.Data, s.Interface()) // Append the interface value of the pointer to the previously created targetType type.
	}

	return res, nil
}
