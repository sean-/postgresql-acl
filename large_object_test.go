package acl_test

import (
	"reflect"
	"testing"

	acl "github.com/sean-/postgresql-acl"
)

func TestLargeObjectString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want acl.LargeObject
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: acl.LargeObject{
				ACL: acl.ACL{
					Role: "foo",
				},
			},
		},
		{
			name: "all without grant",
			in:   "foo=rw",
			out:  "foo=rw",
			want: acl.LargeObject{
				ACL: acl.ACL{
					Role:       "foo",
					Privileges: acl.Select | acl.Update,
				},
			},
		},
		{
			name: "all with grant",
			in:   "foo=r*w*",
			out:  "foo=r*w*",
			want: acl.LargeObject{
				ACL: acl.ACL{
					Role:         "foo",
					Privileges:   acl.Select | acl.Update,
					GrantOptions: acl.Select | acl.Update,
				},
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=r*w*/bar",
			out:  "foo=r*w*/bar",
			want: acl.LargeObject{
				ACL: acl.ACL{
					Role:         "foo",
					GrantedBy:    "bar",
					Privileges:   acl.Select | acl.Update,
					GrantOptions: acl.Select | acl.Update,
				},
			},
		},
		{
			name: "all mixed grant1",
			in:   "foo=rw*",
			out:  "foo=rw*",
			want: acl.LargeObject{
				ACL: acl.ACL{
					Role:         "foo",
					Privileges:   acl.Select | acl.Update,
					GrantOptions: acl.Update,
				},
			},
		},
		{
			name: "all mixed grant2",
			in:   "foo=r*w",
			out:  "foo=r*w",
			want: acl.LargeObject{
				ACL: acl.ACL{
					Role:         "foo",
					Privileges:   acl.Select | acl.Update,
					GrantOptions: acl.Select,
				},
			},
		},
		{
			name: "public all",
			in:   "=r*w*",
			out:  "=r*w*",
			want: acl.LargeObject{
				ACL: acl.ACL{
					Role:         "",
					Privileges:   acl.Select | acl.Update,
					GrantOptions: acl.Select | acl.Update,
				},
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: acl.LargeObject{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: acl.LargeObject{},
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

			got, err := acl.NewLargeObject(aclItem)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse large object ACL %+q: %v", test.in, err)
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
