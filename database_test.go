package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestDatabaseString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Database
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Database{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=CTc",
			out:  "foo=CTc",
			want: pgacl.Database{
				Role:      "foo",
				Create:    true,
				Temporary: true,
				Connect:   true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=C*T*c*",
			out:  "foo=C*T*c*",
			want: pgacl.Database{
				Role:           "foo",
				Connect:        true,
				ConnectGrant:   true,
				Create:         true,
				CreateGrant:    true,
				Temporary:      true,
				TemporaryGrant: true,
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=C*T*c*/bar",
			out:  "foo=C*T*c*/bar",
			want: pgacl.Database{
				Role:           "foo",
				GrantedBy:      "bar",
				Connect:        true,
				ConnectGrant:   true,
				Create:         true,
				CreateGrant:    true,
				Temporary:      true,
				TemporaryGrant: true,
			},
		},
		{
			name: "public all",
			in:   "=c",
			out:  "=c",
			want: pgacl.Database{
				Role:    "",
				Connect: true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Database{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Database{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewDatabase(test.in)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse database ACL %+q: %v", test.in, err)
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
