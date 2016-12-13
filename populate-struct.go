package dbmdl

import (
	"database/sql"
	"errors"
	"reflect"
)

// PopulateStruct loads data from a database a populates the struct
func PopulateStruct(dlct string, t string, s interface{}, limit uint64, whereClauses map[string]interface{}) error {
	d, ok := dialects[dlct]
	if !ok {
		return errors.New("[dbmdl] Failed to register struct; dialect " + dlct + " unknown!")
	}

	ref := reflect.TypeOf(s).Elem()

	if _, exists := tables[ref]; !exists {
		return errors.New("[dbmdl] Type " + ref.Name() + " is not a known type!")
	}

	var fields []string

	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i) // Get the field at index i
		if field.Tag.Get("dbmdl") == "" {
			continue
		}

		fields = append(fields, field.Name)
	}

	// Query
	q := d.FetchFields(t, limit, whereClauses, fields)
	c := make(chan *sql.Rows)
	query(c, q...)

	// Wait for the channel to return rows
	r := <-c
	for r.Next() {
		// TODO: Finish me
	}

	return nil
}
