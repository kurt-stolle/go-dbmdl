package dbmdl

import "database/sql"

// Query is the structure used to handle queries by an external system
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
// Example
//  for {
//     q := <- dbmdl.QueryChannel();
//     sqlDB.Query(q);
//  }
func QueryChannel() chan *Query {
	return chOut
}

// query is our internal query function
func query(res chan *sql.Rows, args ...interface{}) {
	q := new(Query)
	q.String = (args[0]).(string)
	q.Arguments = args[1:]

	if res != nil {
		q.Result = res
	}

	chOut <- q
}
