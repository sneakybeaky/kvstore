package scenario_test

import (
	"kvstore/testing/scenario"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWellFormedEntries(t *testing.T) {

	tests := map[string]struct {
		input string
		want  []scenario.Item
	}{
		"Single GET with expected value": {input: "GET k v\n",
			want: []scenario.Item{
				{Typ: scenario.ItemGet, Val: "GET"},
				{Typ: scenario.ItemKey, Val: " k"},
				{Typ: scenario.ItemValue, Val: " v"},
			},
		},
		"Single PUT": {input: "PUT k v\n",
			want: []scenario.Item{
				{Typ: scenario.ItemPut, Val: "PUT"},
				{Typ: scenario.ItemKey, Val: " k"},
				{Typ: scenario.ItemValue, Val: " v"},
			},
		},
		"Single GET with not found expected": {input: "GET k NOT_FOUND\n",
			want: []scenario.Item{
				{Typ: scenario.ItemGet, Val: "GET"},
				{Typ: scenario.ItemKey, Val: " k"},
				{Typ: scenario.ItemNotFound, Val: " NOT_FOUND"},
			},
		},
		"Multiple entries": {input: "GET k1 NOT_FOUND\nPUT k2 value\nGET k3 value2\n",
			want: []scenario.Item{
				{Typ: scenario.ItemGet, Val: "GET"},
				{Typ: scenario.ItemKey, Val: " k1"},
				{Typ: scenario.ItemNotFound, Val: " NOT_FOUND"},
				{Typ: scenario.ItemPut, Val: "\nPUT"},
				{Typ: scenario.ItemKey, Val: " k2"},
				{Typ: scenario.ItemValue, Val: " value"},
				{Typ: scenario.ItemGet, Val: "\nGET"},
				{Typ: scenario.ItemKey, Val: " k3"},
				{Typ: scenario.ItemValue, Val: " value2"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			l := scenario.NewLexer(t.Name(), test.input)

			var got []scenario.Item

			for {

				i := l.NextItem()

				if i.Typ == scenario.ItemEOF {
					break
				}

				got = append(got, i)

			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf(" mismatch (-want +got):\n%s", diff)
			}

		})
	}

}
