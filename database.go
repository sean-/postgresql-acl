package pgacl

import (
	"bytes"
	"fmt"
	"strings"
)

// Database models the privileges of a database aclitem
type Database struct {
	Role           string
	GrantedBy      string
	Create         bool
	CreateGrant    bool
	Connect        bool
	ConnectGrant   bool
	Temporary      bool
	TemporaryGrant bool
}

const numDatabaseOpts = 6

// NewDatabase parses a PostgreSQL ACL string for a database and returns a Database
// object
func NewDatabase(aclStr string) (Database, error) {
	acl := Database{}
	idx := strings.IndexByte(aclStr, '=')
	if idx == -1 {
		return Database{}, fmt.Errorf("invalid aclStr format: %+q", aclStr)
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
		case 'T':
			acl.Temporary = true
			if withGrant() {
				acl.TemporaryGrant = true
			}
		case 'c':
			acl.Connect = true
			if withGrant() {
				acl.ConnectGrant = true
			}
		case '/':
			if i+1 <= aclLen {
				acl.GrantedBy = aclStr[i+1:]
			}
			break SCAN
		default:
			return Database{}, fmt.Errorf("invalid byte %c in database ACL at %d: %+q", aclStr[i], i, aclStr)
		}
	}

	return acl, nil
}

// String creates a PostgreSQL native output for the ACLs that apply to a
// database.
func (s Database) String() string {
	b := new(bytes.Buffer)
	b.Grow(len(s.Role) + numDatabaseOpts + 1)

	fmt.Fprint(b, s.Role, "=")

	if s.Create {
		fmt.Fprint(b, "C")
		if s.CreateGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.Temporary {
		fmt.Fprint(b, "T")
		if s.TemporaryGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.Connect {
		fmt.Fprint(b, "c")
		if s.ConnectGrant {
			fmt.Fprint(b, "*")
		}
	}

	if s.GrantedBy != "" {
		fmt.Fprint(b, "/", s.GrantedBy)
	}

	return b.String()
}
