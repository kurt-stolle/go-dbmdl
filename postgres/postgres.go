package postgres

import (
	"github.com/kurt-stolle/go-dbmdl"
)

func init() {
	// Set-up the dialeect
	d := new(dbmdl.Dialect)

	// Register for later use in other appliances
	dbmdl.RegisterDialect("postgres", d)
}
