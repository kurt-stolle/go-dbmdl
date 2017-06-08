# DBMDL

A library for modelling databases according to a Go `struct`. Intentionally made lightweight, for this library is not meant to replace SQL languages in your project.

## Basic usage

Structs can be externed with `dbmdl` tags. This tag must always start with the datatype (e.g. `char(15)`) or an `extern` field.

#### Datatype

When a datatype is provided, DBMDL can modify a linked table in the database. In most cases, the programmer would enter a datatype.
The amount of available datatypes depends on the implementation of the SQL driver that is used.

Syntax: `DataType, Parameter1, Parameter2`

##### Parameters

- `primary key`: Indicates that the field is a primary key in the database. There may be multiple primary keys.
- `default X`: Specifies a default value in the database. X indicates some default value.
- `not null`: The value may not be `NULL`. Note that it is preferred to use this, rather than appending `not null` to the database.
- `omit`: Omits the value by default when performing a Population. Useful for columns with optional fields or expensive fields that don't always need to be loaded

#### Select field
The `select` field is used only for loading a struct. It adds a field to the selection statement.

Syntax: `select FieldName`

#### Extern field

The `extern` field is used only for loading a struct. It loads data from an specified table using a JOIN-clause. This is used for when structs need to be linked with data in other tables.

Syntax: `extern <ExternField> from <Table> on <Condition>`. 

Where:
- `ExternField` is a field name of a table that is not our the struct's table
- `Table` is the name of said table
- `Condition` is a joining condition, e.g. `table_one.LocalField2=table_two.ExternField2`

The join type is always `INNER`.

## Example

```go
// Package database
package database

import (
	"github.com/kurt-stolle/go-dbmdl/postgres"
	"database/sql"
)

var Dialect = new(postgres.Dialect)

func Open() *sql.DB {
	return sql.Open(...)	
}	

// Package models
package models

import "github.com/kurt-stolle/go-dbmdl"

type User struct {
    UserID      string          `dbmdl:"uuid, primary key"`
    FirstName   sql.NullString  `dbmdl:"varchar(50)"`
    LastName    string          `dbmdl:"varchar(50), default 'Undefined'"`
    Password    string          `dbmdl:"varchar(255), not null, omit"`
    
    NonDBMDLValue int
    
    privateValue bool
}

var UserModel = dbmdl.NewModel("users", reflect.TypeOf((*User)(nil)), database.Dialect, database.Open)

type Keys struct {
    KeyID           string `dbmdl:"uuid, primary key"`
    LinkedUserID    string `dbmdl:"uuid references users(UserID) on delete cascade, not null"`
    UserLastName    string `dbmdl:"extern users.LastName from users on keys.LinkedUserID=users.UserID"`
}

var KeysModel = dbmdl.NewModel("keys", reflect.TypeOf((*Key)(nil)), database.Dialect, database.Open)
```

## Documentation

See [_godoc_](https://godoc.org/github.com/kurt-stolle/go-dbmdl).
