package postgres

import (
	"strings"

	"github.com/kurt-stolle/go-dbmdl"
)

func init() {
	// Set-up the dialeect
	d := new(dbmdl.Dialect)
	d.CreateTable = func(n string, f ...string) string {
		return `
			CREATE TABLE IF NOT EXISTS ` + n + ` (` + strings.Join(f, ", ") + `);
		`
	}

	// Register for later use in other appliances
	dbmdl.RegisterDialect("postgres", d)
}
