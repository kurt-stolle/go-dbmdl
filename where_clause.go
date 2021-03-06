package dbmdl

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

// WhereClause is a struct for determining the WHERE selectors in the SQL query, for these often vary
type WhereClause struct {
	Clauses []string
	Format  string

	values []interface{}
}

// Values returns a WHERE clause's parameter values
func (w *WhereClause) Values() []interface{} {
	return w.values
}

// String returns a WHERE clause string
func (w *WhereClause) String() string {
	if w.Format != "" {
		return w.FormattedString(w.Format)
	}

	if len(w.Clauses) < 1 {
		return ``
	}

	return `WHERE ` + strings.Join(w.Clauses, " AND ")
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
// Use GetPlaceholder to get a placeholder. This is NOT automatically added to the clause!
// Example: w.AddValuedClause("ID="+w.GetPlaceholder(0), "2015");
func (w *WhereClause) AddValuedClause(clause string, value interface{}) (clauseIndex int, valueIndex int) {
	w.Clauses = append(w.Clauses, clause)
	w.values = append(w.values, value)

	clauseIndex = len(w.Clauses) - 1
	valueIndex = len(w.values) - 1

	return clauseIndex, valueIndex
}

// GetPlaceholder returns a placeholder whose index corresponds to the current amount of entries in values
// In dbmdl, a placeholders are postgres-style: $1, $2, $3, ..., $N
func (w *WhereClause) GetPlaceholder(offset int) string {
	var n = len(w.values) + 1 + offset

	return "$" + strconv.Itoa(n)
}
