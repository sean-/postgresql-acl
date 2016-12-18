# pgacl

## `pgacl` Library

`pgacl` parses
[PostgreSQL's ACL syntax](https://www.postgresql.org/docs/current/static/sql-grant.html#SQL-GRANT-NOTES)
and returns a usable structure.  Library documentation is available at
[https://godoc.org/github.com/sean-/pgacl](https://godoc.org/github.com/sean-/pgacl).


```go
package main

import (
	"fmt"

	"github.com/sean-/pgacl"
)

func main() {
	acl := pgacl.Schema{
		Role:  "foo",
		Usage: true,
	}
	fmt.Printf("ACL String: %s\n", acl.String())

	acl, err := pgacl.NewSchema("foo=C*U")
	if err != nil {
		fmt.Errorf("Bad: %v", err)
	}

	fmt.Printf("ACL Struct: %#v\n", acl)
}
```

```text
ACL String: foo=U
ACL Struct: pgacl.Schema{Role:"foo", Create:true, CreateGrant:true, Usage:true, UsageGrant:false}
```

## Supported PostgreSQL `aclitem` Types

- column permissions
- database
- domain
- foreign data wrappers
- foreign server
- function
- language
- large object
- schema
- sequences
- table
- table space
- type

## Notes

The output from `String()` should match the ordering of characters in `aclitem`
however not all types have been matched with PostgreSQL (yet).

The target of each of these ACLs (e.g. schema name, table name, etc) is not
contained within PostgreSQLs `aclitem` and it is expected this value is managed
elsewhere in your object model.

Arrays of `aclitem` are supposed to be iterated over by the caller, like:

```go
const schema = "public"
var name, owner string
var acls []string
err := conn.QueryRow("SELECT n.nspname, pg_catalog.pg_get_userbyid(n.nspowner), COALESCE(n.nspacl, '{}'::aclitem[])::TEXT[] FROM pg_catalog.pg_namespace n WHERE n.nspname = $1", schema).Scan(&name, &owner, pq.Array(&acls))
if err == nil {
    for _, acl := range acls {
        acl, err = pgacl.NewSchema(acl)
        if err != nil {
            return err
        }
        // ...
    }
}
```
