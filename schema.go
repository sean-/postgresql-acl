package acl

import "fmt"

// Schema models the privileges of a schema aclitem
type Schema struct {
	ACL
}

// NewSchema parses an ACL object and returns a Schema object.
func NewSchema(acl ACL) (Schema, error) {
	if !validRights(acl, validSchemaPrivs) {
		return Schema{}, fmt.Errorf("invalid flags set for schema (%+q), only %+q allowed", permString(acl.Privileges, acl.GrantOptions), validSchemaPrivs)
	}

	return Schema{ACL: acl}, nil
}
