package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// ForeignDataWrapper models the privileges of a function aclitem
type ForeignDataWrapper struct {
	Role       string
	GrantedBy  string
	Usage      bool
	UsageGrant bool
}

const numForeignDataWrapperOpts = 2

// NewForeignDataWrapper parses a PostgreSQL ACL string for a function and returns a ForeignDataWrapper
// object
func NewForeignDataWrapper(aclStr string) (ForeignDataWrapper, error) {
	acl := ForeignDataWrapper{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return ForeignDataWrapper{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
		case 'U':
			acl.Usage = true
			if withGrant() {
				acl.UsageGrant = true
			}
		case '/':
			if i+1 <= aclLen {
				acl.GrantedBy = aclStr[i+1:]
			}
			break SCAN
		default:
			return ForeignDataWrapper{}, fmt.Errorf("invalid byte %c in function ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// function.
func (s ForeignDataWrapper) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(s.Role) + numForeignDataWrapperOpts + 1)

	fmt.Fprint(b, s.Role, "=")

	if s.Usage {
		fmt.Fprint(b, "U")
		if s.UsageGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.GrantedBy != "" {
		fmt.Fprint(b, "/", s.GrantedBy)
	}

	return b.String()
}
