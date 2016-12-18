package acl_test

import (
	"reflect"
	"testing"

	acl "github.com/sean-/postgresql-acl"
)

func TestSchemaString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want acl.Schema
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: acl.Schema{
				ACL: acl.ACL{
					Role: "foo",
				},
			},
		},
		{
			name: "all without grant",
			in:   "foo=UC",
			out:  "foo=UC",
			want: acl.Schema{
				ACL: acl.ACL{
					Role:       "foo",
					Privileges: acl.Create | acl.Usage,
				},
			},
		},
		{
			name: "all with grant",
			in:   "foo=U*C*",
			out:  "foo=U*C*",
			want: acl.Schema{
				ACL: acl.ACL{
					Role:         "foo",
					Privileges:   acl.Create | acl.Usage,
					GrantOptions: acl.Create | acl.Usage,
				},
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=U*C*/bar",
			out:  "foo=U*C*/bar",
			want: acl.Schema{
				ACL: acl.ACL{
					Role:         "foo",
					GrantedBy:    "bar",
					Privileges:   acl.Create | acl.Usage,
					GrantOptions: acl.Create | acl.Usage,
				},
			},
		},
		{
			name: "all mixed grant1",
			in:   "foo=U*C",
			out:  "foo=U*C",
			want: acl.Schema{
				ACL: acl.ACL{
					Role:         "foo",
					Privileges:   acl.Create | acl.Usage,
					GrantOptions: acl.Usage,
				},
			},
		},
		{
			name: "all mixed grant2",
			in:   "foo=UC*",
			out:  "foo=UC*",
			want: acl.Schema{
				ACL: acl.ACL{
					Role:         "foo",
					Privileges:   acl.Create | acl.Usage,
					GrantOptions: acl.Create,
				},
			},
		},
		{
			name: "public all",
			in:   "=U*C*",
			out:  "=U*C*",
			want: acl.Schema{
				ACL: acl.ACL{
					Role:         "",
					Privileges:   acl.Create | acl.Usage,
					GrantOptions: acl.Create | acl.Usage,
				},
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: acl.Schema{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: acl.Schema{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			aclItem, err := acl.Parse(test.in)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse ACLItem %+q: %v", test.in, err)
			}

			if err == nil && test.fail {
				t.Fatalf("expected failure")
			}

			if test.fail && err != nil {
				return
			}

			got, err := acl.NewSchema(aclItem)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse schema ACL %+q: %v", test.in, err)
			}

			if out := test.want.String(); out != test.out {
				t.Fatalf("want %+q got %+q", test.out, out)
			}

			if !reflect.DeepEqual(test.want, got) {
				t.Fatalf("bad: expected %v to equal %v", test.want, got)
			}
		})
	}
}
