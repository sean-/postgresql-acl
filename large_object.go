package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// LargeObject models the privileges of a large object aclitem
type LargeObject struct {
	Role        string
	GrantedBy   string
	Update      bool
	UpdateGrant bool
	Select      bool
	SelectGrant bool
}

const numLargeObjectOpts = 4

// NewLargeObject parses a PostgreSQL ACL string for a large object and returns
// a LargeObject object
func NewLargeObject(aclStr string) (LargeObject, error) {
	acl := LargeObject{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return LargeObject{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
		case '/':
			if i+1 <= aclLen {
				acl.GrantedBy = aclStr[i+1:]
			}
			break SCAN
		default:
			return LargeObject{}, fmt.Errorf("invalid byte %c in large object ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// large object.
func (s LargeObject) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(s.Role) + numLargeObjectOpts + 1)

	fmt.Fprint(b, s.Role, "=")

	if s.Select {
		fmt.Fprint(b, "r")
		if s.SelectGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.Update {
		fmt.Fprint(b, "w")
		if s.UpdateGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.GrantedBy != "" {
		fmt.Fprint(b, "/", s.GrantedBy)
	}

	return b.String()
}
