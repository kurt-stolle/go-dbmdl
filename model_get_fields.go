package dbmdl

func (m *Model) GetFields() ([]string, FromSpecifier) {
	var fields []string
	var clause = new(FromClause)

	// Root tale
	clause.Table = m.TableName

	// Selection loop
FieldLoop:
	for i := 0; i < m.Type.NumField(); i++ {
		field := m.Type.Field(i) // Get the field at index i
		if field.Tag.Get("dbmdl") == "" {
			continue
		}

		params := getTagParameters(field)
		key := params[0]
		params = params[1:]

		if res := regExtern.FindAllString(key, -1); len(res) > 0 {
			// External key
			var extFieldName = res[0]
			var extTableName = res[1]
			var extJoinCondition = res[3]

			// Create a new leaf for the from clause
			clause.AddLeafs(&FromLeaf{
				Table:     extTableName,
				Condition: extJoinCondition,
			})

			// Add the field to the list, but prepend the table name
			fields = append(fields, extTableName+"."+extFieldName)
		} else {
			//	Data type definition
			for _, s := range params {
				if s == omit {
					continue FieldLoop // Omits the field when searching for fields to load in selections
				}
			}

			fields = append(fields, field.Name)
		}
	}

	return fields, clause
}
