package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// Function models the privileges of a function aclitem
type Function struct {
	Role         string
	GrantedBy    string
	Execute      bool
	ExecuteGrant bool
}

const numFunctionOpts = 2

// NewFunction parses a PostgreSQL ACL string for a function and returns a Function
// object
func NewFunction(aclStr string) (Function, error) {
	acl := Function{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return Function{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
		case 'X':
			acl.Execute = true
			if withGrant() {
				acl.ExecuteGrant = true
			}
		case '/':
			if i+1 <= aclLen {
				acl.GrantedBy = aclStr[i+1:]
			}
			break SCAN
		default:
			return Function{}, fmt.Errorf("invalid byte %c in function ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// function.
func (s Function) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(s.Role) + numFunctionOpts + 1)

	fmt.Fprint(b, s.Role, "=")

	if s.Execute {
		fmt.Fprint(b, "X")
		if s.ExecuteGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.GrantedBy != "" {
		fmt.Fprint(b, "/", s.GrantedBy)
	}

	return b.String()
}
