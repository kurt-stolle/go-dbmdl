package dbmdl

func (m *Model) GetFields() ([]*FieldMapping, FromSpecifier) {
	var fields []*FieldMapping
	var clause = new(FromClause)

	// Root tale
	clause.table = m.TableName

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

		if res := regExtern.FindStringSubmatch(key); len(res) > 0 {
			// External key
			var extFieldName = res[1]
			var extTableName = res[2]
			var extJoinCondition = res[3]
			var extJoinType string

			if resJoinType := regExternJoin.FindStringSubmatch(extJoinCondition); len(res) > 0 {
				extJoinCondition = resJoinType[1]
				extJoinType = resJoinType[2]
			} else {
				extJoinType = "INNER"
			}

			// Create a new leaf for the from clause
			clause.AddLeaf(&FromLeaf{
				Table:     extTableName,
				JoinType:  extJoinType,
				Condition: extJoinCondition,
			})

			// Add the field to the list, but prepend the table name
			fields = append(fields, &FieldMapping{
				Link:   field.Name,
				Clause: extFieldName,
			})
		} else if res := regSelect.FindStringSubmatch(key); len(res) > 0 {
			fields = append(fields, &FieldMapping{
				Link:   field.Name,
				Clause: res[1],
			})
		} else {
			//	Data type definition
			for _, s := range params {
				if s == omit {
					continue FieldLoop // Omits the field when searching for fields to load in selections
				}
			}

			fields = append(fields, &FieldMapping{
				Link: field.Name,
			})
		}
	}

	return fields, clause
}
