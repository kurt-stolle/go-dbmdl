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

	d.CreateTable = func(n string, f []string) string {
		var query bytes.Buffer
		query.WriteString(`CREATE TABLE IF NOT EXISTS ` + n + ` ();`)

		for _, dt := range f {
			query.WriteString(`DO $$
	BEGIN
		BEGIN
			ALTER TABLE ` + n + ` ADD COLUMN ` + dt + `;
		EXCEPTION
			WHEN duplicate_column THEN RAISE NOTICE 'column ` + dt + ` already exists in ` + n + `.';
		END;
	END;
$$;`)
		}

		return query.String()
	}

	d.SetPrimaryKey = func(n string, f []string) string {
		return `DO $$
	BEGIN
		if not exists (select constraint_name
	  	from information_schema.constraint_column_usage
			where table_name = '` + n + `' and constraint_name = '` + n + `_pkey') then
	    	execute 'ALTER TABLE ` + n + ` ADD PRIMARY KEY (` + strings.Join(f, ",") + `)';
	  end if;
	END;
$$;`
	}

	d.SetDefaultValues = func(n string, v map[string]string) string {
		var q []string
		for c, d := range v {
			q = append(q, `UPDATE `+n+` SET `+c+`=`+d+` WHERE `+c+`=NULL;
				ALTER TABLE ONLY `+n+` ALTER COLUMN `+c+` SET DEFAULT `+d+`;`)
		}

		return strings.Join(q, "\n")
	}

	d.FetchFields = func(tableName string, fields []string, p *dbmdl.Pagination, w *dbmdl.WhereClause) (string, []interface{}) {
		var query bytes.Buffer
		var args []interface{}

		// Basic query
		query.WriteString(`SELECT `)
		query.WriteString(strings.Join(fields, ", "))
		query.WriteString(` FROM `)
		query.WriteString(tableName)

		// Where clauses
		if w != nil {
			query.WriteString(` ` + w.String() + ` `)
			args = w.Values
		}

		// Pagination
		if p != nil {
			query.WriteString(` ` + p.String() + ` `)
		}

		// Result
		return query.String(), args
	}

	d.Insert = func(tableName string, fieldsValues map[string]interface{}) (string, []interface{}) {
		var args = []interface{}{}
		var query bytes.Buffer

		query.WriteString(`INSERT INTO `)
		query.WriteString(tableName)
		query.WriteString(` (`)

		var bufInsert string
		var bufValues string

		var i = 0
		for f, v := range fieldsValues {
			i++

			bufInsert += f

			bufValues += "$"
			bufValues += strconv.Itoa(i)

			if i < len(fieldsValues) {
				bufInsert += ","
				bufValues += ","
			}

			args = append(args, v)
		}

		query.WriteString(bufInsert)
		query.WriteString(`) VALUES (`)
		query.WriteString(bufValues)
		query.WriteString(`)`)

		return query.String(), args
	}

	d.Update = func(tableName string, fieldsValues map[string]interface{}, w *dbmdl.WhereClause) (string, []interface{}) {
		var args = []interface{}{}
		var query bytes.Buffer

		args = append(args, w.Values...)

		query.WriteString(`UPDATE `)
		query.WriteString(tableName)
		query.WriteString(` SET `)

		var i = len(w.Values)
		for f, v := range fieldsValues {
			i++

			query.WriteString(f)
			query.WriteString(`=$`)
			query.WriteString(strconv.Itoa(i))

			args = append(args, v)
		}

		query.WriteString(w.String())

		return query.String(), args
	}

	d.GetPlaceholder = func(offset int) string {
		return "$" + strconv.Itoa(offset)
	}

	// Register for later use in other appliances
	dbmdl.RegisterDialect("postgres", d)
}
