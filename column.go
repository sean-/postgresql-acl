package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// Column models the privileges of a column aclitem
type Column struct {
	Role            string
	GrantedBy       string
	Insert          bool
	InsertGrant     bool
	References      bool
	ReferencesGrant bool
	Select          bool
	SelectGrant     bool
	Update          bool
	UpdateGrant     bool
}

const numColumnOpts = 8

// NewColumn parses a PostgreSQL ACL string for a column and returns a Column
// object
func NewColumn(aclStr string) (Column, error) {
	acl := Column{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return Column{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
	}

	acl.Role = aclStr[:idx]

	aclLen := len(aclStr)
	var i int
	withGrant := func() bool {
		if i+1 >= aclLen {
			return false
		}

		if aclStr[i+1] == '*' {
			i++
			return true
		}

		return false
	}

SCAN:
	for i = idx + 1; i < aclLen; i++ {
		switch aclStr[i] {
		case 'a':
			acl.Insert = true
			if withGrant() {
				acl.InsertGrant = true
			}
		case 'r':
			acl.Select = true
			if withGrant() {
				acl.SelectGrant = true
			}
		case 'w':
			acl.Update = true
			if withGrant() {
				acl.UpdateGrant = true
			}
		case 'x':
			acl.References = true
			if withGrant() {
				acl.ReferencesGrant = true
			}
		case '/':
			if i+1 <= aclLen {
				acl.GrantedBy = aclStr[i+1:]
			}
			break SCAN
		default:
			return Column{}, fmt.Errorf("invalid byte %c in column ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// column.
func (t Column) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(t.Role) + numColumnOpts + 1)

	fmt.Fprint(b, t.Role, "=")

	if t.Insert {
		fmt.Fprint(b, "a")
		if t.InsertGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.Select {
		fmt.Fprint(b, "r")
		if t.SelectGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.Update {
		fmt.Fprint(b, "w")
		if t.UpdateGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.References {
		fmt.Fprint(b, "x")
		if t.ReferencesGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.GrantedBy != "" {
		fmt.Fprint(b, "/", t.GrantedBy)
	}

	return b.String()
}
