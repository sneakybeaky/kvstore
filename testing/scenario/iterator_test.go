package scenario_test

import (
	"kvstore/testing/scenario"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWellFormedGet(t *testing.T) {

	want := []*scenario.Scenario{
		{
			Operation: scenario.Get,
			Key:       "key",
			Value:     "wantedValue",
		},
	}
	var got []*scenario.Scenario

	mni := memoryNextItem{items: []scenario.Item{
		{
			Typ: scenario.ItemGet,
			Val: "GET",
		},
		{
			Typ: scenario.ItemKey,
			Val: "key",
		},
		{
			Typ: scenario.ItemValue,
			Val: "wantedValue",
		},
		{
			Typ: scenario.ItemEOF,
		},
	}}

	for s, err := range scenario.IterateScenarios(mni.nextItem) {
		if err != nil {
			t.Fatal(err)
		}

		got = append(got, s)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf(" mismatch (-want +got):\n%s", diff)
	}

}

type memoryNextItem struct {
	items []scenario.Item
	pos   int
}

func (m *memoryNextItem) nextItem() scenario.Item {
	defer func() { m.pos++ }()
	return m.items[m.pos]

}
