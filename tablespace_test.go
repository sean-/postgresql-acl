package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestTablespaceString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Tablespace
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Tablespace{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=C",
			out:  "foo=C",
			want: pgacl.Tablespace{
				Role:   "foo",
				Create: true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=C*",
			out:  "foo=C*",
			want: pgacl.Tablespace{
				Role:        "foo",
				Create:      true,
				CreateGrant: true,
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=C*/bar",
			out:  "foo=C*/bar",
			want: pgacl.Tablespace{
				Role:        "foo",
				GrantedBy:   "bar",
				Create:      true,
				CreateGrant: true,
			},
		},
		{
			name: "public all",
			in:   "=C",
			out:  "=C",
			want: pgacl.Tablespace{
				Role:   "",
				Create: true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Tablespace{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Tablespace{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewTablespace(test.in)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse tablespace ACL %+q: %v", test.in, err)
			}

			if err == nil && test.fail {
				t.Fatalf("expected failure")
			}

			if test.fail && err != nil {
				return
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
