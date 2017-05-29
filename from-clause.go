package dbmdl

// FromLeaf specifies a conditional joining of tables, i.e. a JOIN clause
type FromLeaf struct {
	Table     string
	JoinType  string // JoinType specifies "left", "right" or "inner"
	Condition string
}

// FromClause implements FromSpecifier, used to determine which table to select things from
type FromClause struct {
	Table string
	Leafs []FromLeaf
}

// String created a string of the form `FROM <table_name> JOIN <specification>`
func (fc *FromClause) String() string {
	var str = `FROM ` + fc.Table

	for _, l := range fc.Leafs {
		if l.JoinType == "" {
			str += " INNER "
		} else {
			str += " " + l.JoinType
		}
		str += " JOIN " + l.Table + " ON " + l.Condition
	}

	return str
}