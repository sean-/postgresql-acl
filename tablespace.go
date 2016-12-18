package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// Tablespace models the privileges of a tablespace aclitem
type Tablespace struct {
	Role        string
	GrantedBy   string
	Create      bool
	CreateGrant bool
}

const numTablespaceOpts = 2

// NewTablespace parses a PostgreSQL ACL string for a tablespace and returns a Tablespace
// object
func NewTablespace(aclStr string) (Tablespace, error) {
	acl := Tablespace{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return Tablespace{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
		case 'C':
			acl.Create = true
			if withGrant() {
				acl.CreateGrant = true
			}
		case '/':
			if i+1 <= aclLen {
				acl.GrantedBy = aclStr[i+1:]
			}
			break SCAN
		default:
			return Tablespace{}, fmt.Errorf("invalid byte %c in tablespace ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// tablespace.
func (s Tablespace) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(s.Role) + numTablespaceOpts + 1)

	fmt.Fprint(b, s.Role, "=")

	if s.Create {
		fmt.Fprint(b, "C")
		if s.CreateGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.GrantedBy != "" {
		fmt.Fprint(b, "/", s.GrantedBy)
	}

	return b.String()
}
