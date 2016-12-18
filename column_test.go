package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestColumnString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Column
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Column{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=arwx",
			out:  "foo=arwx",
			want: pgacl.Column{
				Role:       "foo",
				Insert:     true,
				References: true,
				Select:     true,
				Update:     true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=a*r*w*x*",
			out:  "foo=a*r*w*x*",
			want: pgacl.Column{
				Role:            "foo",
				Insert:          true,
				InsertGrant:     true,
				References:      true,
				ReferencesGrant: true,
				Select:          true,
				SelectGrant:     true,
				Update:          true,
				UpdateGrant:     true,
			},
		},
		{
			name: "all with grant and by",
			in:   "foo=a*r*w*x*/bar",
			out:  "foo=a*r*w*x*/bar",
			want: pgacl.Column{
				Role:            "foo",
				GrantedBy:       "bar",
				Insert:          true,
				InsertGrant:     true,
				References:      true,
				ReferencesGrant: true,
				Select:          true,
				SelectGrant:     true,
				Update:          true,
				UpdateGrant:     true,
			},
		},
		{
			name: "public all",
			in:   "=r",
			out:  "=r",
			want: pgacl.Column{
				Role:   "",
				Select: true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Column{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Column{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewColumn(test.in)
			if err != nil && !test.fail {
				t.Fatalf("unable to parse table ACL %+q: %v", test.in, err)
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
