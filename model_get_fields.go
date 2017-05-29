package dbmdl

import (
	"log"
	"strings"
)

func (m *Model) GetFields() ([]string, FromSpecifier) {
	var fields []string
	var clause = new(FromClause)

	clause.Table = m.TableName

	for i := 0; i < m.Type.NumField(); i++ {
		field := m.Type.Field(i) // Get the field at index i
		if field.Tag.Get("dbmdl") == "" {
			continue
		}

		params := getTagParameters(field)

		if params[0] == "extern" {
			if len(params) != 2 {
				log.Fatalf("Field %s in struct %s has invalid extern tag, 2 parameters and no more must be provided!", field.Name, m.Type.Name())
			}

			d := strings.Split(params[1], " ")
			l := len(d)
			if l < 5 || l > 6 || d[1] != "from" || d[3] != "on" {
				log.Fatalf("Field %s in struct %s has invalid extern tag; the 2nd parameter must be of the form: <ExternField> from <Table> on <Condition> [JoinType]", field.Name, m.Type.Name())
			}

			if l == 6 {
				clause.Leafs = append(clause.Leafs, FromLeaf{
					Table:     d[2],
					Condition: d[4],
					JoinType:  d[5],
				})
			} else {
				clause.Leafs = append(clause.Leafs, FromLeaf{
					Table:     d[2],
					Condition: d[4],
					JoinType:  "inner",
				})
			}
		} else {
			for _, tag := range params {
				if tag == omit {
					continue
				}
			}

			fields = append(fields, field.Name)
		}

	}

	return fields, clause
}
