package dbmdl

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
)

// createTables will register the struct in the database
func createTables(ref reflect.Type) error {
	var t, ok = tables[ref]
	if !ok {
		return errors.New("[dbmdl] Type not in tables map: " + ref.Name())
	}

	// Build fields list
	var fields []string
	var primaryKeys []string

	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)          // Get the field at index i
		tag := getTagParameters(field) // Find the datatype from the dbmdl tag

		if len(tag) <= 0 || tag[0] == "" {
			continue
		}

		fields = append(fields, field.Name+" "+tag[0])

		for _, v := range tag {
			switch v {
			case "primary key":
				primaryKeys = append(primaryKeys, field.Name)
			}
		}
	}

	// Query
	if len(primaryKeys) <= 0 {
		log.Fatal("[dbmdl] Struct " + ref.Name() + " has no primary key")
	}

	// A query
	var q []interface{}

	// Channel magic
	c1 := make(chan *sql.Rows)                // Create a new channel;
	q = t.dialect.CreateTable(t.name, fields) //  Make the table query
	query(c1, q...)                           // Perform query
	<-c1                                      // Wait for query to finish

	c2 := make(chan *sql.Rows)                       // Make another channel
	q = t.dialect.SetPrimaryKey(t.name, primaryKeys) // Build primary key query
	query(c2, q...)                                  // Execute query
	<-c2                                             // Wait for query to finish

	return nil
}

//
