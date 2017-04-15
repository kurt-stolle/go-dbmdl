# DBMDL

A library for modelling databases according to a Go `struct`. Intentionally made lightweight, for this library is not meant to replace SQL languages in your project.

## Usage

### Registering a dialect that queries will be constructed in

You can choose the dialect for your queries by importing a package that defines one, or calling `dbmdl.RegisterDialect(name string, d *dbmdl.Dialect)`. In the following example, we register the Postgres dialect.

```go
import (
  ...
  _ "github.com/kurt-stolle/go-dbmdl/postgres" //
)
```

The "\_" character indicates that we only import this package for the side effects, i.e. the `init()` function which registers the dialect.

### Registering the structs that can be used by dbmdl

To use a struct in dbmdl, it must be registered first.

```go
type MyModel struct {
  Key      int    `dbmdl:"serial, primary key"`
  Value    string `dbmdl:"varchar(100)"`
}

var conn *sql.DB = postgres.Connect()
if err := dbmdl.RegisterStruct(conn, "postgres", "project_models", (*MyModel)(nil)); err != nil { // (*MyModel)(nil) allows us to pass the type only so that we can use it in reflection
  panic(err);
}
```

Evident from the example above, fields must have a `dbmdl` tag to save them in the database. The struct must also have at least one `primary key`.

### The `dbmdl` tag

The `dbmdl` struct field tag must always start with the (database) datatype, e.g. `char(15)` or an `extern` field.

#### Datatype

When a datatype is provided, DBMDL can modify a linked table in the database. In most cases, the programmer would enter a datatype.
The amount of available datatypes depends on the implementation of the SQL driver that is used.

#### Extern field

The `extern` field is used only for loading a struct. It loads data from an specified table using a JOIN-clause. This is used for when structs need to be linked with data in other tables.

The syntax of the `extern` field is as follows: `extern Field at Table from LocalForeignKeyField`.

#### Additional parameters

- `primary key`: Indicates that the field is a primary key in the database. There may be multiple primary keys.
- `default X`: Specifies a default value in the database. X indicates some default value.
- `not null`: The value may not be `NULL`. Note that it is preferred to use this, rather than appending `not null` to the database.
- `omit`: Omits the value by default when performing a Population. Useful for columns with optional fields or expensive fields that don't always need to be loaded

Parameters are separated by a comma sign.

#### Examples

```go
type Model struct {
  Index     int     `dbmdl:"serial, primary key, not null"`
  ValueOne  string  `dbmdl:"varchar(50), primary key, not null"`
  ValueTwo  string  `dbmdl:"varchar(90), primary key, not null, default 'Something'"`
}
```

## Documentation

See [_godoc_](https://godoc.org/github.com/kurt-stolle/go-dbmdl).
