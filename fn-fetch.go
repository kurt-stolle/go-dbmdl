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
func Fetch(dlct string, t string, sRef interface{}, FWP ...interface{}) (*Result, error) {
	d, ok := dialects[dlct]
	if !ok {
		return nil, errors.New("[dbmdl] Failed to populate struct; dialect " + dlct + " unknown!")
	}

	ref := reflect.TypeOf(sRef).Elem()

	if _, exists := tables[ref]; !exists {
		return nil, errors.New("[dbmdl] Type " + ref.Name() + " is not a known type!")
	}

	var fields []string
	var where *WhereClause
	var pag *Pagination

	for i := 0; i < len(FWP); i++ {
		switch v := FWP[i].(type) {
		case []string:
			fields = v
		case *Pagination:
			pag = v
		case *WhereClause:
			where = v
		default:
			panic("Interface type not supported")
		}
	}

	if pag == nil {
		pag = &Pagination{1, 1}
	}

	if where != nil {
		where.Dialect = d
	}

	if len(fields) < 1 { // If we did not supply and fields to be selected, select all fields
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

	if len(fields) < 1 {
		return nil, errors.New("We have nothing to select because all fields are ommited by dbmdl")
	}

	// Query
	q := d.FetchFields(t, fields, pag, where)

	c := make(chan *sql.Rows)
	query(c, q...)

	// Results
	var res = &Result{}
	res.Pagination = pag

	// Create dummy variables that we can scan the results of the query into
	var dummyVariables []reflect.Value // Slice to hold the values scanned from the *sql.Rows result
	var dummyVariablesAddresses []interface{}
	for _, name := range fields {
		f, found := ref.FieldByName(name)
		if !found {
			return nil, errors.New("Field not found in array: " + name)
		}

		var nwtyp = reflect.New(f.Type)

		dummyVariables = append(dummyVariables, nwtyp.Elem()) // Create a pointer to type at f
		dummyVariablesAddresses = append(dummyVariablesAddresses, nwtyp.Interface())
	}

	// Wait for the channel to return rows
	r := <-c
	defer r.Close()
	defer close(c)

	// Iterate over rows found
	for r.Next() {
		var s = reflect.New(ref) // Create a new pointer to an empty struct of type ref

		r.Scan(dummyVariablesAddresses...) // Scan into the slice we populated with dummy variables earlier

		for i, v := range dummyVariables {
			s.Elem().FieldByName(fields[i]).Set(v) // Set values in our new struct
		}

		res.Data = append(res.Data, s.Interface()) // Append the interface value of the pointer to the previously created ref type.
	}

	return res, nil
}
