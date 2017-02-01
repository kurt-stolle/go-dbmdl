package dbmdl

import "strconv"

// Pagination is used to prevent dbmdl from selecting far too many entries
type Pagination struct {
	Page          uint
	AmountPerPage uint
}

// String returns a LIMIT and OFFSET clause
func (p *Pagination) String() string {
	return `LIMIT ` + strconv.FormatUint(uint64(p.AmountPerPage), 10) + ` OFFSET ` + strconv.FormatUint(uint64((p.Page-1)*p.AmountPerPage), 10)
}

// NewPagination returns a Pagination pointer
func NewPagination(page, amount uint) *Pagination {
	return &Pagination{page, amount}
}
