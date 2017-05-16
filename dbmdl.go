// Package dbmdl is used for modelling a database according to a golang struct.
package dbmdl

import (
    "database/sql"
    "errors"
    "reflect"
)

// Constants
const (
    omit = "omit"
)

// Globals
var (
    // Errors
    ErrNotFound       = sql.ErrNoRows
    ErrStructNotFound = errors.New("dbmdl: The struct provided is not registered with DBMDL")
    ErrNoDialect      = errors.New("dbmdl: WhereClause has no dialect set")
    ErrNoPointer      = errors.New("dbmdl: Target is not a pointer")
    ErrUnknownType    = errors.New("dbmdl: Unknown type requested")

    // Field properties
)

// A where selector selects which rows must be selected. Implemented by WhereClause
type WhereSelector interface {
    String() string
    Values() []interface{}
}

// Translator is an interface implemented by language dialects, such as github.com/kurt-stolle/go-dbmdl/postgres.Dialect
type Translator interface {
    CreateTable(tableName string) string
    AddField(tableName, field, def string) string
    SetPrimaryKeys(tableName string, fields []string) string
    SetDefaultValue(n, field, def string) string
    SetNotNull(n, field string) string
    FetchFields(tableName string, fields []string, p *Pagination, w WhereSelector) (string, []interface{})
    Insert(tableName string, fieldsValues map[string]interface{}) (string, []interface{})
    Update(tableName string, fieldsValues map[string]interface{}, w WhereSelector) (string, []interface{})
    Count(tableName string, w WhereSelector) (string, []interface{})
    GetPlaceholder(i int) string
}


// Model is a modeller tied to a struct
type Model struct {
    TableName   string
    Type        reflect.Type
    Dialect     Translator
    GetDatabase func() *sql.DB
}

// NewModeller creates a new Model for a certain database and reflection type
// Modellers should be saved and re-used
func NewModel(tableName string, reflectionType reflect.Type, dialect Translator, getDatabaseFunc func() *sql.DB) *Model {
    for reflectionType.Kind() == reflect.Ptr {
        reflectionType = reflectionType.Elem()
    }

    mdl := new(Model)
    mdl.TableName = tableName
    mdl.Type = reflectionType
    mdl.Dialect = dialect
    mdl.GetDatabase = getDatabaseFunc

    return mdl
}
