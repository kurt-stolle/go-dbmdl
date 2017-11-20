package dbmdl

// FromLeaf specifies a conditional joining of tables, i.e. a JOIN clause
type FromLeaf struct {
	Table     string
	JoinType  string
	Condition string
}

// FromClause implements FromSpecifier, used to determine which table to select things from
type FromClause struct {
	table string
	Leafs []*FromLeaf
}

// AddLeaf appends a new FromLeaf
func (fc *FromClause) AddLeaf(l *FromLeaf) {
	for _, c := range fc.Leafs {
		if l.Table == c.Table {
			return // Only allow 1 leaf per table
		}
	}
	fc.Leafs = append(fc.Leafs, l)
}

// String created a string of the form `FROM <table_name> JOIN <specification>`
func (fc *FromClause) String() string {
	var str = `FROM ` + fc.table

	for _, l := range fc.Leafs {
		str += " " + l.JoinType + " JOIN " + l.Table + " ON " + l.Condition
	}

	return str
}

// GetTable returns the root table
func (fc *FromClause) GetTable() string {
	return fc.table
}
