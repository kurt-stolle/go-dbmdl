package dbmdl

import "reflect"

// tagImpliesNotStored returns true if the dbmdl tag's value indicates that the value is not in the database
func fieldNotInDatabase(f reflect.StructField) bool {
	switch f.Tag.Get("dbmdl") {
	case "-": fallthrough
	case "select": fallthrough
	case "extern": fallthrough
	case "": return true
	}

	return false
}
