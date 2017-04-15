package dbmdl

import (
	"database/sql"
	"errors"
	"log"
)

// GenerateTable creates a table for the struct in the database
func GenerateTable(db *sql.DB, reference interface{}) error {
	var ref = getReflectType(reference)
	var t, ok = tables[ref]
	if !ok {
		return errors.New("Type not in tables map: " + ref.Name())
	}

	// Build fields list
	var fields []string
	var primaryKeys []string
	var notNull []string
	var defaults = make(map[string]string)

	// Iterate over fields
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)          // Get the field at index i
		tag := getTagParameters(field) // Find the datatype from the dbmdl tag

		if len(tag) <= 0 || tag[0] == "" {
			continue
		}

		// Is this an extern?
		if regExtern.MatchString(tag[0]) {
			continue
		}

		// First find the special tags
		for _, v := range tag {
			if i := regDefault.FindStringIndex(v); i != nil {
				defaults[field.Name] = v[(i[0] + 8):] // Move 8 spaces to the right from 'default ' to capture the type
			} else if v == "primary key" {
				primaryKeys = append(primaryKeys, field.Name)
			} else if v == "not null" {
				notNull = append(notNull, field.Name)
			}
		}

		// Add the definition to the list
		fields = append(fields, field.Name+" "+tag[0])
	}

	// Query
	if len(primaryKeys) <= 0 {
		log.Fatal("dbmdl: Struct " + ref.Name() + " has no primary key")
	}

	// Start generating the database.
	// This is done synchronously to prevent any problems with database RW
	var q string

	q = t.dialect.CreateTable(t.name)
	if _, err := db.Exec(q); err != nil {
		log.Fatal("dbmdl: Failed to create table ", t.name, "\nQuery:", q, "\nError:", err)
	}

	for _, field := range fields {
		q = t.dialect.AddField(t.name, field)
		if _, err := db.Exec(q); err != nil {
			log.Fatal("dbmdl: Failed to add column ", field, "\nQuery:", q, "\nError:", err)
		}
	}

	q = t.dialect.SetPrimaryKeys(t.name, primaryKeys)
	if _, err := db.Exec(q); err != nil {
		log.Fatal("dbmdl: Failed to set primary keys for ", t.name, "\nQuery:", q, "\nError:", err)
	}

	for field, def := range defaults {
		q = t.dialect.SetDefaultValue(t.name, field, def)
		if _, err := db.Exec(q); err != nil {
			log.Fatal("dbmdl: Failed to set default value for column ", field, " to ", def, "\nQuery:", q, "\nError:", err)
		}
	}

	for _, field := range notNull {
		q = t.dialect.SetNotNull(t.name, field)
		if _, err := db.Exec(q); err != nil {
			log.Fatal("dbmdl: Failed to add not null constraint to column", field, "\nQuery:", q, "\nError:", err)
		}
	}

	return nil
}

//
