package postgres

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/kurt-stolle/go-dbmdl"
)

func init() {
	// Set-up the dialeect
	d := new(dbmdl.Dialect)
	d.CreateTable = func(n string, f []string) []interface{} {
		var args []interface{}
		var query bytes.Buffer
		query.WriteString(`CREATE TABLE IF NOT EXISTS ` + n + ` ();`)

		for _, dt := range f {
			query.WriteString(`ALTER TABLE ` + n + ` ADD COLUMN ` + dt + `;`)
		}

		// Build arg list
		args[0] = query.String()
		return args
	}
	d.FetchFields = func(n string, limit uint64, w map[string]interface{}, f []string) []interface{} {
		var args []interface{}
		args[0] = "" // Save this spot for later

		var query bytes.Buffer
		query.WriteString(`SELECT ` + strings.Join(f, ", ") + ` FROM ` + n)

		var whereClauses []string
		for key, val := range w {
			whereClauses = append(whereClauses, key+`=$`+strconv.Itoa(len(whereClauses)+1))
			args = append(args, val)
		}

		if len(whereClauses) > 0 {
			query.WriteString(` WHERE `)
			query.WriteString(strings.Join(whereClauses, " AND "))
		}

		if limit > 0 {
			query.WriteString(` LIMIT ` + strconv.FormatUint(limit, 10))
		}

		args[0] = query.String() // Put query string in reserved spot

		return args
	}

	// Register for later use in other appliances
	dbmdl.RegisterDialect("postgres", d)
}
