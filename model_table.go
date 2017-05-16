package dbmdl

import (
	"log"
)

// CreateTable creates a table for the struct in the database
func (m *Model) CreateTable() error {
	// Build fields list
	var fields []([2]string)
	var primaryKeys []string
	var notNull []string
	var defaults = make(map[string]string)

	// Iterate over fields
	for i := 0; i < m.Type.NumField(); i++ {
		field := m.Type.Field(i)          // Get the field at index i
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
		fields = append(fields, [2]string{field.Name,tag[0]})
	}

	// Query
	if len(primaryKeys) <= 0 {
		log.Fatal("dbmdl: Struct " + m.Type.Name() + " has no primary key")
	}

	// Start generating the database.
	// This is done synchronously to prevent any problems with database RW
	var q string
	var dl = m.Dialect
	var db = m.GetDatabase()
	var tableName = m.TableName

	q = m.Dialect.CreateTable(tableName)
	if _, err := db.Exec(q); err != nil {
		log.Fatal("dbmdl: Failed to create table ", tableName, "\nQuery:", q, "\nError:", err)
	}

	for _, field := range fields {
		q = dl.AddField(tableName, field[0], field[1])
		if _, err := db.Exec(q); err != nil {
			log.Fatal("dbmdl: Failed to add column ", field, "\nQuery:", q, "\nError:", err)
		}
	}

	q = dl.SetPrimaryKeys(tableName, primaryKeys)
	if _, err := db.Exec(q); err != nil {
		log.Fatal("dbmdl: Failed to set primary keys for ", tableName, "\nQuery:", q, "\nError:", err)
	}

	for field, def := range defaults {
		q = dl.SetDefaultValue(tableName, field, def)
		if _, err := db.Exec(q); err != nil {
			log.Fatal("dbmdl: Failed to set default value for column ", field, " to ", def, "\nQuery:", q, "\nError:", err)
		}
	}

	for _, field := range notNull {
		q = dl.SetNotNull(tableName, field)
		if _, err := db.Exec(q); err != nil {
			log.Fatal("dbmdl: Failed to add not null constraint to column", field, "\nQuery:", q, "\nError:", err)
		}
	}

	return nil
}

//
