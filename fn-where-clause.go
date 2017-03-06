package dbmdl

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

// WhereClause is a struct for determining the WHERE selectors in the SQL query, for these often vary
type WhereClause struct {
	Values  []interface{}
	Clauses []string
	Dialect *Dialect // Dialect usually needn't be set by the programmer, it is implicitly found by the relevant dbmdl functions
	Format  string
	Order   []string
}

// helper function for the order
func orderString(or []string) string {
	if len(or) < 1 {
		return ""
	}

	return " ORDER BY " + strings.Join(or, ",")
}

// NewWhereClause returns a where clause with a dialect
// Accepts:
// (string) dialect name
// (*Dialect) dialect
// (reflect.Type) type
// (-other-) -> reflect.Type
func NewWhereClause(ifc interface{}) *WhereClause {
	w := new(WhereClause)
	switch v := ifc.(type) {
	case string:
		d, ok := dialects[v]
		if !ok {
			panic(ErrNoDialect)
		}

		w.Dialect = d
	case *Dialect:
		w.Dialect = v
	default:
		d, ok := tables[getReflectType(v)]
		if !ok {
			panic(ErrStructNotFound)
		}

		w.Dialect = d.dialect
	}

	return w
}

// String returns a WHERE clause string
func (w *WhereClause) String() string {
	if w.Format != "" {
		return w.FormattedString(w.Format)
	}

	if len(w.Clauses) < 1 {
		return ``
	}

	return `WHERE ` + strings.Join(w.Clauses, " AND ") + orderString(w.Order)
}

// FormattedString strings a WHERE clause string
// In contrast to String(), this function also accepts a format used to concatenate the Clauses into a string
// Format is simply a number for the index, e.g. 0 AND 1 AND (2 OR 3)
// This is useful when not all clauses are connected by AND
func (w *WhereClause) FormattedString(f string) string {
	var reg = regexp.MustCompile(`/\d+/g`)
	var buf bytes.Buffer
	var cursor = 0

	for _, ran := range reg.FindAllStringIndex(f, -1) {
		buf.WriteString(f[cursor : ran[0]-1])

		i, err := strconv.Atoi(f[ran[0]:ran[1]])
		if err != nil {
			panic(err)
		} else if i >= len(w.Clauses) {
			panic("Format has indeces that are out of range. Attempted format: " + f)
		}

		buf.WriteString(w.Clauses[i])
	}

	return buf.String() + orderString(w.Order)
}

// AddClause a clause and possibly a value
func (w *WhereClause) AddClause(clause string) (clauseIndex int) {
	w.Clauses = append(w.Clauses, clause)

	clauseIndex = len(w.Clauses) - 1

	return clauseIndex
}

// AddValuedClause adds a clause with a value
// Use GetPlaceholder to get a placeholder. This is NOT automatically added to the clause!
// Example: w.AddValuedClause("ID="+w.GetPlaceholder(0), "2015");
func (w *WhereClause) AddValuedClause(clause string, value interface{}) (clauseIndex int, valueIndex int) {
	w.Clauses = append(w.Clauses, clause)
	w.Values = append(w.Values, value)

	clauseIndex = len(w.Clauses) - 1
	valueIndex = len(w.Values) - 1

	return clauseIndex, valueIndex
}

// OrderDesc orders the set descending
func (w *WhereClause) OrderDesc(column string) {
	w.Order = append(w.Order, column+" DESC")
}

// OrderAsc orders the set ascending
func (w *WhereClause) OrderAsc(column string) {
	w.Order = append(w.Order, column+" ASC")
}

// GetPlaceholder returns a placeholder whose index corresponds to the current amount of entries in Values
func (w *WhereClause) GetPlaceholder(offset int) string {
	return w.Dialect.GetPlaceholder(len(w.Values) + 1 + offset)
}
