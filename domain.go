package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// Domain models the privileges of a domain aclitem
type Domain struct {
	Role       string
	GrantedBy  string
	Usage      bool
	UsageGrant bool
}

const numDomainOpts = 2

// NewDomain parses a PostgreSQL ACL string for a domain and returns a Domain
// object
func NewDomain(aclStr string) (Domain, error) {
	acl := Domain{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return Domain{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
			return Domain{}, fmt.Errorf("invalid byte %c in domain ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// domain.
func (s Domain) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(s.Role) + numDomainOpts + 1)

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
