package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestDomainString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Domain
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Domain{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=U",
			out:  "foo=U",
			want: pgacl.Domain{
				Role:  "foo",
				Usage: true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=U*",
			out:  "foo=U*",
			want: pgacl.Domain{
				Role:       "foo",
				Usage:      true,
				UsageGrant: true,
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=U*/bar",
			out:  "foo=U*/bar",
			want: pgacl.Domain{
				Role:       "foo",
				GrantedBy:  "bar",
				Usage:      true,
				UsageGrant: true,
			},
		},
		{
			name: "all mixed grant1",
			in:   "foo=U*",
			out:  "foo=U*",
			want: pgacl.Domain{
				Role:       "foo",
				Usage:      true,
				UsageGrant: true,
			},
		},
		{
			name: "all mixed grant2",
			in:   "foo=U",
			out:  "foo=U",
			want: pgacl.Domain{
				Role:  "foo",
				Usage: true,
			},
		},
		{
			name: "public all",
			in:   "=U*",
			out:  "=U*",
			want: pgacl.Domain{
				Role:       "",
				Usage:      true,
				UsageGrant: true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Domain{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Domain{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewDomain(test.in)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse domain ACL %+q: %v", test.in, err)
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
