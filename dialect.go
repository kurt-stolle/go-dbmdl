package dbmdl

// Dialect is a struct that stores the querying methods
type Dialect struct {
	CreateTable func(tableName string, fields ...string) string
}

// Create the channel
var ch chan string

func init() {
	ch = make(chan string)
}

// QueryChannel returns a channel for executing in own implementation
// Example:
//  for {
//     q := <- dbmdl.QueryChannel();
//     sqlDB.Query(q);
//  }
func QueryChannel() chan string {
	return ch
}

// query is our internal query function
func query(q string) {
	ch <- q
}
