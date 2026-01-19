package scenario_test

import (
	"kvstore/testing/scenario"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWellFormedGet(t *testing.T) {
	l := scenario.NewLexer(t.Name(), "GET k v\n")

	var got []scenario.Item

	want := []scenario.Item{
		{Typ: scenario.ItemGet, Val: "GET"},
		{Typ: scenario.ItemKey, Val: " k"},
		{Typ: scenario.ItemValue, Val: " v"},
	}

	for {

		i := l.NextItem()

		if i.Typ == scenario.ItemEOF {
			break
		}

		got = append(got, i)

	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(" mismatch (-want +got):\n%s", diff)
	}

}
