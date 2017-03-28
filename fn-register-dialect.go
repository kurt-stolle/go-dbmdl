package dbmdl

import "log"

// Dialect is a struct that stores the querying methods
type Dialect struct {
	CreateTable     func(tableName string) string
	AddField        func(tableName string, field string) string
	SetPrimaryKeys  func(tableName string, fields []string) string
	SetDefaultValue func(n string, field string, def string) string
	SetNotNull      func(n string, field string) string
	FetchFields     func(tableName string, fields []string, p *Pagination, w *WhereClause) (string, []interface{})
	Insert          func(tableName string, fieldsValues map[string]interface{}) (string, []interface{})
	Update          func(tableName string, fieldsValues map[string]interface{}, w *WhereClause) (string, []interface{})
	Count           func(tableName string, w *WhereClause) (string, []interface{})
	GetPlaceholder  func(i int) string
}

// RegisterDialect will add a dialect so that it can be used later
func RegisterDialect(d string, strct *Dialect) error {
	log.Println("Registered dialect: " + d)

	dialects[d] = strct

	return nil
}
