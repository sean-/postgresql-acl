package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestLargeObjectString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.LargeObject
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.LargeObject{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=rw",
			out:  "foo=rw",
			want: pgacl.LargeObject{
				Role:   "foo",
				Select: true,
				Update: true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=r*w*",
			out:  "foo=r*w*",
			want: pgacl.LargeObject{
				Role:        "foo",
				Select:      true,
				SelectGrant: true,
				Update:      true,
				UpdateGrant: true,
			},
		},
		{
			name: "all with grant by role",
			in:   "foo=r*w*/bar",
			out:  "foo=r*w*/bar",
			want: pgacl.LargeObject{
				Role:        "foo",
				GrantedBy:   "bar",
				Select:      true,
				SelectGrant: true,
				Update:      true,
				UpdateGrant: true,
			},
		},
		{
			name: "public all",
			in:   "=r",
			out:  "=r",
			want: pgacl.LargeObject{
				Role:   "",
				Select: true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.LargeObject{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.LargeObject{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewLargeObject(test.in)
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
