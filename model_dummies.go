package dbmdl

import (
	"errors"
	"reflect"
)

// GetDummies returns matching dummy variables' addresses for scanning the result of a query into
func (m *Model) GetDummies(fields ...string) ([]interface{}, error) {
	// Slice of addresses
	var addr = make([]interface{},len(fields))

	// Iterate over requested fields
	for i,f := range fields {
		tp,ok := m.Type.FieldByName(f)
		if !ok {
			return nil, errors.New("dbmdl: field "+f+" does not exist")
		} else if fieldNotInDatabase(tp) {
			return nil, errors.New("dbmdl: field "+f+" is not a field replicated in a database")
		}

		addr[i] = reflect.New(tp.Type).Interface() // Create a new variable and store the address at the index matching the fields slice
	}

	// Return addresses, no error
	return addr,nil
}

// MapDummies maps a slice of addresses t o a slice of field names with matching indexes, useful after scanning into the addresses
// provided by Model.GetDummies
func MapDummies(fields []string, addr []interface{}) (map[string]interface{}, error) {
	if len(fields) != len(addr) {
		return nil, errors.New("dbmdl: slice dimensions do not agree")
	}

	var mp = make(map[string]interface{})
	for i,f := range fields {
		mp[f] = reflect.ValueOf(addr[i]).Elem().Interface() // Using reflection, find the value pointed to by addr[i] and convert back to interface{}
	}

	return mp, nil
}