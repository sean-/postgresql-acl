package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// Table models the privileges of a table aclitem
type Table struct {
	Role            string
	GrantedBy       string
	Delete          bool
	DeleteGrant     bool
	Insert          bool
	InsertGrant     bool
	References      bool
	ReferencesGrant bool
	Select          bool
	SelectGrant     bool
	Trigger         bool
	TriggerGrant    bool
	Truncate        bool
	TruncateGrant   bool
	Update          bool
	UpdateGrant     bool
}

const numTableOpts = 14

// NewTable parses a PostgreSQL ACL string for a table and returns a Table
// object
func NewTable(aclStr string) (Table, error) {
	acl := Table{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return Table{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
		case 'd':
			acl.Delete = true
			if withGrant() {
				acl.DeleteGrant = true
			}
		case 'D':
			acl.Truncate = true
			if withGrant() {
				acl.TruncateGrant = true
			}
		case 'r':
			acl.Select = true
			if withGrant() {
				acl.SelectGrant = true
			}
		case 't':
			acl.Trigger = true
			if withGrant() {
				acl.TriggerGrant = true
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
			return Table{}, fmt.Errorf("invalid byte %c in table ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// table.
func (t Table) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(t.Role) + numTableOpts + 1)

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

	if t.Delete {
		fmt.Fprint(b, "d")
		if t.DeleteGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.Truncate {
		fmt.Fprint(b, "D")
		if t.TruncateGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.References {
		fmt.Fprint(b, "x")
		if t.ReferencesGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.Trigger {
		fmt.Fprint(b, "t")
		if t.TriggerGrant {
			fmt.Fprint(b, "*")
		}
	}

	if t.GrantedBy != "" {
		fmt.Fprint(b, "/", t.GrantedBy)
	}

	return b.String()
}
