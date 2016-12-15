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
