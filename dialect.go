package dbmdl

import "database/sql"

// Dialect is a struct that stores the querying methods
type Dialect struct {
	CreateTable func(tableName string, fields []string) []interface{}
	FetchFields func(tableName string, limit uint64, whereClauses map[string]interface{}, fields []string) []interface{}
}

// Query handles queries
type Query struct {
	String    string
	Arguments []interface{}
	Result    chan *sql.Rows
}

// Create the channel
var chOut chan *Query

func init() {
	chOut = make(chan *Query)
}

// QueryChannel returns a channel for executing in own implementation
// Example:
//  for {
//     q := <- dbmdl.QueryChannel();
//     sqlDB.Query(q);
//  }
func QueryChannel() chan *Query {
	return chOut
}

// query is our internal query function
func query(res chan *sql.Rows, args ...interface{}) {
	go func() {
		q := new(Query)
		q.String = args[0].(string)
		q.Arguments = args[1:]

		if res != nil {
			q.Result = make(chan *sql.Rows)
		}

		chOut <- q
	}()
}
