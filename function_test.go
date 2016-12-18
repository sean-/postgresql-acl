package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestFunctionString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Function
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Function{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=X",
			out:  "foo=X",
			want: pgacl.Function{
				Role:    "foo",
				Execute: true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=X*",
			out:  "foo=X*",
			want: pgacl.Function{
				Role:         "foo",
				Execute:      true,
				ExecuteGrant: true,
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=X*/bar",
			out:  "foo=X*/bar",
			want: pgacl.Function{
				Role:         "foo",
				GrantedBy:    "bar",
				Execute:      true,
				ExecuteGrant: true,
			},
		},
		{
			name: "all mixed grant1",
			in:   "foo=X*",
			out:  "foo=X*",
			want: pgacl.Function{
				Role:         "foo",
				Execute:      true,
				ExecuteGrant: true,
			},
		},
		{
			name: "all mixed grant2",
			in:   "foo=X",
			out:  "foo=X",
			want: pgacl.Function{
				Role:    "foo",
				Execute: true,
			},
		},
		{
			name: "public all",
			in:   "=X*",
			out:  "=X*",
			want: pgacl.Function{
				Role:         "",
				Execute:      true,
				ExecuteGrant: true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Function{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Function{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewFunction(test.in)
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