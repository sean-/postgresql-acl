package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestSequenceString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Sequence
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Sequence{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=rwU",
			out:  "foo=rwU",
			want: pgacl.Sequence{
				Role:   "foo",
				Select: true,
				Update: true,
				Usage:  true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=r*w*U*",
			out:  "foo=r*w*U*",
			want: pgacl.Sequence{
				Role:        "foo",
				Select:      true,
				SelectGrant: true,
				Update:      true,
				UpdateGrant: true,
				Usage:       true,
				UsageGrant:  true,
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=r*w*U*/bar",
			out:  "foo=r*w*U*/bar",
			want: pgacl.Sequence{
				Role:        "foo",
				GrantedBy:   "bar",
				Select:      true,
				SelectGrant: true,
				Update:      true,
				UpdateGrant: true,
				Usage:       true,
				UsageGrant:  true,
			},
		},
		{
			name: "public all",
			in:   "=rU",
			out:  "=rU",
			want: pgacl.Sequence{
				Role:   "",
				Select: true,
				Usage:  true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Sequence{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Sequence{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewSequence(test.in)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse sequence ACL %+q: %v", test.in, err)
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
