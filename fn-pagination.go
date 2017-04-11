package dbmdl

import (
	"bytes"
	"database/sql"
	"math"
	"strconv"
	"strings"
)

// Sorting returns an array of how a result will be sorted
type Sorting []sortfield

type sortfield struct {
	Field     string
	Ascending bool
}

// Pagination is used to prevent dbmdl from selecting far too many entries
type Pagination struct {
	Page          int
	AmountPerPage int
	Prev          int
	Next          int
	First         int
	Last          int

	Sorting Sorting
}

// String returns a LIMIT and OFFSET clause
func (p *Pagination) String() string {
	return strings.Join([]string{p.SortingString(), p.LimitsString()}, " ") // Use join because we don't know whether either exists
}

// LimitsString returns the limits and offset
func (p *Pagination) LimitsString() string {
	return `LIMIT ` + strconv.Itoa(p.AmountPerPage) + ` OFFSET ` + strconv.Itoa((p.Page-1)*p.AmountPerPage)
}

// SortingString returns the sorting order
func (p *Pagination) SortingString() string {
	if len(p.Sorting) <= 0 {
		return ""
	}
	var nw bytes.Buffer

	nw.WriteString("ORDER BY")
	for _, spec := range p.Sorting {
		nw.WriteByte(' ')
		nw.WriteString(spec.Field)
		nw.WriteByte(' ')
		if spec.Ascending {
			nw.WriteString("ASC")
			continue
		}
		nw.WriteString("DESC")
	}

	return nw.String()
}

// OrderDescending adds a descending order item
func (p *Pagination) OrderDescending(field string) {
	p.Sorting = append(p.Sorting, sortfield{field, false})
}

// OrderAscending adds a ascending order item
func (p *Pagination) OrderAscending(field string) {
	p.Sorting = append(p.Sorting, sortfield{field, false})
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
