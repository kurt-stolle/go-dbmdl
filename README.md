# DBMDL

A library for modelling databases according to a Go `struct`. Intentionally made very lightweight, for this library is not meant to replace SQL languages.

## Usage

### 1\. Registering a dialect that queries will be constructed in

You can choose the dialect for your queries by importing a package that defines one, or calling `dbmdl.RegisterDialect(name string, d *dbmdl.Dialect)`. In the following example, we register the Postgres dialect.

```
import (
  ...
  _ "github.com/kurt-stolle/go-dbmdl/postgres" //
)
```

The "\_" character indicates that we only import this package for the side effects, i.e. the `init()` function which registers the dialect.

### 2\. Setting up a channel for receiving queries built by dbmdl

A channel needs to be created in order to receive the queries constructed by dbmdl.

```
go func() {                        // Create a new goroutine for DBMDL queries
    ch := dbmdl.QueryChannel()       // Initialize the channel
    for {                            // Keep waiting for queries
        q := <-ch                      // Receive a query

    conn, err := postgres.Open()   // Open a connection
        if err != nil {                // Check for errors
            log.Fatal(err)
        }
        defer conn.Close()             // Defer connection close (required)

        rows, err := conn.Query(q.String, q.Arguments...) // Execute the query
        if err != nil {                                   // Check for errors
            log.Fatal("Failed to execute DBMDL query! \nQuery: ", q.String, "\nError: ", err)
        }
    defer rows.Close()     // Defer rows close (required)

        if q.Result != nil {   // If a result return is requested by dbmdl, grant it.
            q.Result <- rows     // Send rows to the result channel
        }
    }
}()
```

### 3\. Registering the structs that can be used by dbmdl.

To use a struct in dbmdl, it must be registered first.

```
type MyModel struct {
  Key      int    `dbmdl:"serial, primary key"`
  Value    string `dbmdl:"varchar(100)"`
}
if err := dbmdl.RegisterStruct("postgres", "project_models", &MyModel{}); err != nil {
  panic(err);
}
```

Evident from the example above, fields must have a `dbmdl` tag in order for them to be saved to the database. The struct must also have at least one `primary key`.

## The `dbmdl` tag

The `dbmdl` tag must always start with the (database) datatype. This means that this datatype is **never** implicit from the Golang datatype.

Other fields that may optionally be added to the tag are:

- `primary key`: Indicates that the field is a primary key in the database. There may be multiple primary keys.
- `default X`: Specifies a default value in the database. X indicates some default value.
- `not null`: The value may not be `NULL`. Note that it is preferred to use this, rather than appending `not null` to the database.
- `omit`: Omits the value by default when performing a Population. Useful for columns with optional fields or expensive fields that don't always need to be loaded

Optional fields are separated from the primary field and other optional field by a "," character.

The following is an example of an elaborate `dbmdl` tag:

```
type Model struct {
  Index     int     `dbmdl:"serial, primary key, not null"`
  ValueOne  string  `dbmdl:"varchar(50), primary key, not null"`
  ValueTwo  string  `dbmdl:"varchar(90), primary key, not null, default 'Something'"`
}
```

## Pagination

...

## WhereClauses

...
