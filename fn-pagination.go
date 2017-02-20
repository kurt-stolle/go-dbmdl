package dbmdl

import (
	"database/sql"
	"math"
	"strconv"
	"strings"
)

// Constants for ordering
const (
	OrderAscending  = " ASC"
	OrderDescending = " DESC"
)

// Pagination is used to prevent dbmdl from selecting far too many entries
type Pagination struct {
	Page          int
	AmountPerPage int
	OrderBy       string
	Prev          int
	Next          int
	First         int
	Last          int
}

// String returns a LIMIT and OFFSET clause
func (p *Pagination) String() string {
	return strings.Join([]string{p.OrderBy, `LIMIT`, strconv.Itoa(p.AmountPerPage), ` OFFSET `, strconv.Itoa((p.Page - 1) * p.AmountPerPage)}, " ")
}

// Order specifies the sorting Order
func (p *Pagination) Order(spec ...string) {
	p.OrderBy = "ORDER BY " + strings.Join(spec, ",")
}

// Load will populate the struct according to a table name and where clause
func (p *Pagination) Load(db *sql.DB, t string, where *WhereClause) error {
	// First is always 1
	p.First = 1

	// Amount of rows that satisfy the where clause
	var count int

	// Perform a query to get the count
	query, args := where.Dialect.Count(t, where)
	if err := db.QueryRow(query, args...).Scan(&count); err != nil {
		return err
	}

	// Fill the pagination table with the rows we selected
	p.Last = int(math.Ceil(float64(count) / float64(p.AmountPerPage)))
	p.Prev = p.Page - 1
	if p.Prev < p.First {
		p.Prev = p.First
	}
	p.Next = p.Page + 1
	if p.Next > p.Last {
		p.Next = p.Last
	}

	return nil
}

// NewPagination returns a Pagination pointer
func NewPagination(page, amount int) *Pagination {
	pag := new(Pagination)
	pag.Page = page
	pag.AmountPerPage = amount
	pag.First = 1

	return pag
}
