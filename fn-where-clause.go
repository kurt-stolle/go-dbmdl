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
	Dialect *Dialect
	Format  string
}

// String returns a WHERE clause string
func (w *WhereClause) String() string {
	if w.Format != "" {
		return w.FormattedString(w.Format)
	}

	str := strings.Join(w.Clauses, " AND ")

	if str == "" {
		return ""
	}
	return " WHERE " + str + " "
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

	return buf.String()
}

// AddClause a clause and possibly a value
func (w *WhereClause) AddClause(clause string) (clauseIndex int) {
	w.Clauses = append(w.Clauses, clause)

	clauseIndex = len(w.Clauses) - 1

	return clauseIndex
}

// AddValuedClause adds a clause with a value
func (w *WhereClause) AddValuedClause(clause string, value interface{}) (clauseIndex int, valueIndex int) {
	w.Clauses = append(w.Clauses, clause)
	w.Values = append(w.Values, value)

	clauseIndex = len(w.Clauses) - 1
	valueIndex = len(w.Values) - 1

	return clauseIndex, valueIndex
}

// GetPlaceholder returns a placeholder whose index corresponds to the current amount of entries in Values
func (w *WhereClause) GetPlaceholder(offset int) string {
	return w.Dialect.GetPlaceholder(len(w.Values) + 1 + offset)
}
