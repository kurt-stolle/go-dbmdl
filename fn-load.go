package dbmdl

import (
	"database/sql"
	"errors"
	"reflect"
)

// Load will load a single struct from the database based on a where clause
// Target is a pointer to a struct
func Load(dlct string, table string, target interface{}, where *WhereClause) error {
	// Check if target is actually a pointer
	ol := reflect.ValueOf(target)
	if ol.Kind() == reflect.Ptr {
		ol = ol.Elem()
	}

	// Check whether the dialect exists
	d, ok := dialects[dlct]
	if !ok {
		return errors.New("[dbmdl] Dialect " + dlct + " unknown")
	}

	// Set the dialect of the WhereClause
	where.Dialect = d

	// Get the fields
	var fields []string
	var ref = reflect.TypeOf(target)
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

	// Query using the same shit as Fetch Fields
	q := d.FetchFields(table, fields, &Pagination{1, 1}, where)
	c := make(chan *sql.Rows)
	query(c, q...)

	// Wait for query to return a result and start scanning
	r := <-c
	for r.Next() {
		var a []interface{} // Slice to hold the values scanned

		for _, name := range fields {
			f, found := ref.FieldByName(name)
			if !found {
				return errors.New("[dbmdl] Field not found in array: " + name)
			}
			a = append(a, reflect.New(f.Type).Interface()) // Populate a by getting the type of each field and getting an infc
		}

		var aAddresses []interface{} // Make an array for the addresses of the values

		for i, aval := range a { // Populate the addresses array
			aAddresses[i] = &aval
		}

		r.Scan(aAddresses...) // Scan into a by pointer reference

		for i, v := range a {
			ol.FieldByName(fields[i]).Set(reflect.ValueOf(v)) // Set values in our new struct
		}
	}

	return nil
}
