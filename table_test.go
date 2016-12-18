package pgacl_test

import (
	"reflect"
	"testing"

	"github.com/sean-/pgacl"
)

func TestTableString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		want pgacl.Table
		fail bool
	}{
		{
			name: "default",
			in:   "foo=",
			out:  "foo=",
			want: pgacl.Table{Role: "foo"},
		},
		{
			name: "all without grant",
			in:   "foo=arwdDxt",
			out:  "foo=arwdDxt",
			want: pgacl.Table{
				Role:       "foo",
				Delete:     true,
				Insert:     true,
				References: true,
				Select:     true,
				Trigger:    true,
				Truncate:   true,
				Update:     true,
			},
		},
		{
			name: "all with grant",
			in:   "foo=a*r*w*d*D*x*t*",
			out:  "foo=a*r*w*d*D*x*t*",
			want: pgacl.Table{
				Role:            "foo",
				Delete:          true,
				DeleteGrant:     true,
				Insert:          true,
				InsertGrant:     true,
				References:      true,
				ReferencesGrant: true,
				Select:          true,
				SelectGrant:     true,
				Trigger:         true,
				TriggerGrant:    true,
				Truncate:        true,
				TruncateGrant:   true,
				Update:          true,
				UpdateGrant:     true,
			},
		},
		{
			name: "all with grant and by",
			in:   "foo=a*r*w*d*D*x*t*/bar",
			out:  "foo=a*r*w*d*D*x*t*/bar",
			want: pgacl.Table{
				Role:            "foo",
				GrantedBy:       "bar",
				Delete:          true,
				DeleteGrant:     true,
				Insert:          true,
				InsertGrant:     true,
				References:      true,
				ReferencesGrant: true,
				Select:          true,
				SelectGrant:     true,
				Trigger:         true,
				TriggerGrant:    true,
				Truncate:        true,
				TruncateGrant:   true,
				Update:          true,
				UpdateGrant:     true,
			},
		},
		{
			name: "public all",
			in:   "=r",
			out:  "=r",
			want: pgacl.Table{
				Role:   "",
				Select: true,
			},
		},
		{
			name: "invalid input1",
			in:   "bar*",
			want: pgacl.Table{},
			fail: true,
		},
		{
			name: "invalid input2",
			in:   "%",
			want: pgacl.Table{},
			fail: true,
		},
	}

	for i, test := range tests {
		if test.name == "" {
			t.Fatalf("test %d needs a name", i)
		}

		t.Run(test.name, func(t *testing.T) {
			got, err := pgacl.NewTable(test.in)
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
