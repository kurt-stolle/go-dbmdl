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

	d.CreateTable = func(n string) string {
		return `CREATE TABLE IF NOT EXISTS ` + n + ` ();`
	}

	d.AddField = func(n, f string) string {
		return `DO $$
	BEGIN
		BEGIN
			ALTER TABLE ` + n + ` ADD COLUMN ` + f + `;
		EXCEPTION
			WHEN duplicate_column THEN RAISE NOTICE 'column ` + f + ` already exists in ` + n + `.';
		END;
	END;
$$;`
	}

	d.SetPrimaryKeys = func(n string, f []string) string {
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

	d.SetDefaultValue = func(n, field, def string) string {
		return `UPDATE ` + n + ` SET ` + field + `=` + def + ` WHERE ` + field + ` IS NULL;
				ALTER TABLE ` + n + ` ALTER COLUMN ` + field + ` SET DEFAULT ` + def
	}

	d.SetNotNull = func(n string, v string) string {
		return `ALTER TABLE ` + n + ` ALTER COLUMN ` + v + ` SET NOT NULL`
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

	d.Count = func(tableName string, w *dbmdl.WhereClause) (string, []interface{}) {
		return ("SELECT COUNT(*) AS rows FROM " + tableName + " " + w.String()), w.Values // Pretty slow, but this is used only for pagination - a less accurate method would not  be sufficient
	}

	d.GetPlaceholder = func(offset int) string {
		return "$" + strconv.Itoa(offset)
	}

	// Register for later use in other appliances
	dbmdl.RegisterDialect("postgres", d)
}
