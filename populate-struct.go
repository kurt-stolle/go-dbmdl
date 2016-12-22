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

// PopulateStruct loads data from a database a populates the struct
// sRef is a pointer to the struct, only used for getting the reflection type
func PopulateStruct(dlct string, t string, sRef interface{}, FWP ...interface{}) (*Result, error) {
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

	if where == nil {
		where = &WhereClause{}
	}

	if len(fields) < 1 { // If we did not supply and fields to be selected, select all fields
		for i := 0; i < ref.NumField(); i++ {
			field := ref.Field(i) // Get the field at index i
			if field.Tag.Get("dbmdl") == "" {
				continue
			}

			for _, tag := range getTagParameters(field) {
				if tag == "omit" {
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
	q := d.FetchFields(t, pag, where, fields)
	c := make(chan *sql.Rows)
	query(c, q...)

	// Results
	var res = &Result{}
	res.Pagination = pag

	// Wait for the channel to return rows
	r := <-c
	for r.Next() {
		var s = reflect.New(ref) // Create a new struct
		var a []interface{}      // Slice to hold the values scanned

		for _, name := range fields {
			f, found := ref.FieldByName(name)
			if !found {
				return nil, errors.New("Field not found in array: " + name)
			}
			a = append(a, reflect.New(f.Type).Interface()) // Populate a by getting the type of each field and getting an infc
		}

		var aAddresses []interface{} // Make an array for the addresses of the values

		for i, aval := range a { // Populate the addresses array
			aAddresses[i] = &aval
		}

		r.Scan(aAddresses...) // Scan into a by pointer reference

		for i, v := range a {
			s.FieldByName(fields[i]).Set(reflect.ValueOf(v)) // Set values in our new struct
		}

		res.Data = append(res.Data, s.Interface()) // Append the interface
	}

	return res, nil
}
