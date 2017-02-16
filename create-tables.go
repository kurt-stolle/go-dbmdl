package dbmdl

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
	"regexp"
)

// createTables will register the struct in the database
func createTables(db *sql.DB, ref reflect.Type) error {
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
	var q string

	q = t.dialect.CreateTable(t.name, fields) //  Make the table query
	if _, err := db.Exec(q); err != nil {
		log.Panic(q)
	}

	q = t.dialect.SetPrimaryKey(t.name, primaryKeys) // Build primary key query
	if _, err := db.Exec(q); err != nil {
		log.Panic(q)
	}

	q = t.dialect.SetDefaultValues(t.name, defaults) // Build default values query
	if _, err := db.Exec(q); err != nil {
		log.Panic(q)
	}

	return nil
}

//
