package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// Sequence models the privileges of a sequence aclitem
type Sequence struct {
	Role        string
	GrantedBy   string
	Update      bool
	UpdateGrant bool
	Select      bool
	SelectGrant bool
	Usage       bool
	UsageGrant  bool
}

const numSequenceOpts = 6

// NewSequence parses a PostgreSQL ACL string for a sequence and returns a Sequence
// object
func NewSequence(aclStr string) (Sequence, error) {
	acl := Sequence{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return Sequence{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
			return Sequence{}, fmt.Errorf("invalid byte %c in sequence ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// sequence.
func (s Sequence) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(s.Role) + numSequenceOpts + 1)

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
