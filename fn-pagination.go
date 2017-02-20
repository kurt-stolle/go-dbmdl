package dbmdl

import (
	"database/sql"
	"math"
	"strconv"
)

// Pagination is used to prevent dbmdl from selecting far too many entries
type Pagination struct {
	Page          uint
	AmountPerPage uint
	Prev          uint
	Next          uint
	First         uint
	Last          uint
}

// String returns a LIMIT and OFFSET clause
func (p *Pagination) String() string {
	return `LIMIT ` + strconv.FormatUint(uint64(p.AmountPerPage), 10) + ` OFFSET ` + strconv.FormatUint(uint64((p.Page-1)*p.AmountPerPage), 10)
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
	p.Last = uint(math.Ceil(float64(count) / float64(p.AmountPerPage)))
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
func NewPagination(page, amount uint) *Pagination {
	pag := new(Pagination)
	pag.Page = page
	pag.AmountPerPage = amount
	pag.First = 1

	return pag
}