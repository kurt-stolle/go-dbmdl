package dbmdl

import "log"

// Dialect is a struct that stores the querying methods
type Dialect struct {
	CreateTable      func(tableName string, fields []string) []interface{}
	SetPrimaryKey    func(tableName string, fields []string) []interface{}
	SetDefaultValues func(n string, v map[string]string) []interface{}
	FetchFields      func(tableName string, fields []string, p *Pagination, w *WhereClause) []interface{}
	Insert           func(tableName string, fieldsValues map[string]interface{}) []interface{}
	Update           func(tableName string, fieldsValues map[string]interface{}, w *WhereClause) []interface{}
	GetPlaceholder   func(i int) string
}

// RegisterDialect will add a dialect so that it can be used later
func RegisterDialect(d string, strct *Dialect) error {
	log.Println("[dbmdl] Registered dialect: " + d)

	dialects[d] = strct

	return nil
}
