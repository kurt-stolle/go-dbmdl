package dbmdl

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"regexp"
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
	var defaults = make(map[string]string)

	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)          // Get the field at index i
		tag := getTagParameters(field) // Find the datatype from the dbmdl tag

		if len(tag) <= 0 || tag[0] == "" {
			continue
		}

		fields = append(fields, field.Name+" "+tag[0])

		regDefault := regexp.MustCompile("default .+")

		for _, v := range tag {
			if v == "primary key" {
				primaryKeys = append(primaryKeys, field.Name)
			} else if i := regDefault.FindStringIndex(v); i != nil {
				defaults[field.Name] = v[(i[0] + 8):] // Move 8 spaces to the right from 'default ' to capture the type
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

	c3 := make(chan *sql.Rows)                       // Make another channel
	q = t.dialect.SetDefaultValues(t.name, defaults) // Build primary key query
	query(c3, q...)                                  // Execute query
	<-c3                                             // Wait for query to finish

	return nil
}

//
