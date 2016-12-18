package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestSchemaString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Schema
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Schema{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=UC",
			out:  "foo=UC",
			want: pgacl.Schema{
				Role:   "foo",
				Create: true,
				Usage:  true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=U*C*",
			out:  "foo=U*C*",
			want: pgacl.Schema{
				Role:        "foo",
				Create:      true,
				CreateGrant: true,
				Usage:       true,
				UsageGrant:  true,
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=U*C*/bar",
			out:  "foo=U*C*/bar",
			want: pgacl.Schema{
				Role:        "foo",
				GrantedBy:   "bar",
				Create:      true,
				CreateGrant: true,
				Usage:       true,
				UsageGrant:  true,
			},
		},
		{
			name: "all mixed grant1",
			in:   "foo=U*C",
			out:  "foo=U*C",
			want: pgacl.Schema{
				Role:       "foo",
				Create:     true,
				Usage:      true,
				UsageGrant: true,
			},
		},
		{
			name: "all mixed grant2",
			in:   "foo=UC*",
			out:  "foo=UC*",
			want: pgacl.Schema{
				Role:        "foo",
				Create:      true,
				CreateGrant: true,
				Usage:       true,
			},
		},
		{
			name: "public all",
			in:   "=U*C*",
			out:  "=U*C*",
			want: pgacl.Schema{
				Role:        "",
				Create:      true,
				CreateGrant: true,
				Usage:       true,
				UsageGrant:  true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Schema{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Schema{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewSchema(test.in)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse schema ACL %+q: %v", test.in, err)
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
