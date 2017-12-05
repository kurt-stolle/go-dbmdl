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

		rawTag := field.Tag.Get("dbmdl")
		if rawTag == "extern" {
			// External key
			var extFieldName = field.Tag.Get("dbmdl_field")
			var extTableName = field.Tag.Get("dbmdl_table")
			var extJoinCondition = field.Tag.Get("dbmdl_condition")
			var extJoinType = field.Tag.Get("dbmdl_join")

			if extFieldName == "" {
				panic("dbmdl: Field '" + field.Name + "' in '" + m.Type.Name() + "' recognized as extern, but misses dbmdl_field tag")
			} else if extTableName == "" {
				panic("dbmdl: Field '" + field.Name + "' in '" + m.Type.Name() + "' recognized as extern, but misses dbmdl_table tag")
			} else if extJoinCondition == "" {
				panic("dbmdl: Field '" + field.Name + "' in '" + m.Type.Name() + "' recognized as extern, but misses dbmdl_condition tag")
			}

			if extJoinType == "" {
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
		} else if rawTag == "select" {
			var cls = field.Tag.Get("dbmdl_field")
			if cls == "" {
				panic("dbmdl: Field '" + field.Name + "' in '" + m.Type.Name() + "' recognized as select, but misses dbmdl_field tag")
			}
			fields = append(fields, &FieldMapping{
				Link:   field.Name,
				Clause: cls,
			})
		} else {
			params := getTagParameters(rawTag)[1:]
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
