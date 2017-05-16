package dbmdl

import (
	"reflect"
)

// Save will add to the database or update an existing resource if a nonzero where clause is provided
func (m *Model) Save(target interface{}, where WhereSelector, fields ...string) error {
	// If there are no fields provided, select every field without an omit tag
	if len(fields) < 1 {
		fields = m.GetFields()
	}

	// Build fieldsValues
	var fieldsValues = make(map[string]interface{})
	var val = reflect.ValueOf(target).Elem()

	for _, f := range fields {
		fieldsValues[f] = val.FieldByName(f).Interface()
	}

	// Handle query
	var q string
	var a []interface{}
	if where == nil { // If the where clause is empty, INSERT:
		q, a = m.Dialect.Insert(m.TableName, fieldsValues) // Build query
	} else { // If the where clause is not empty, UPDATE:
		q, a = m.Dialect.Update(m.TableName, fieldsValues, where) // Build query
	}


	// Wait for response and close channel
	if _, err := m.GetDatabase().Exec(q, a...); err != nil {
		return err
	}

	return nil
}
