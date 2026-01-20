package scenario

import (
	"errors"
	"iter"
)

type Operation int

const (
	Get Operation = iota
	Put
	GetNotFound
)

type Scenario struct {
	Operation Operation
	Key       string
	Value     string
}

func IterateScenarios(
	nextItem func() Item,
) iter.Seq2[*Scenario, error] {

	return func(yield func(*Scenario, error) bool) {
		for {

			item := nextItem()

			if item.Typ == ItemEOF {
				return
			}

			if item.Typ == ItemError {
				yield(nil, errors.New(item.Val))
			}

			_ = item.Val

			item = nextItem()

			if item.Typ == ItemError {
				yield(nil, errors.New(item.Val))
			}

			key := item.Val

			item = nextItem()

			if item.Typ == ItemError {
				yield(nil, errors.New(item.Val))
			}

			value := item.Val

			if !yield(&Scenario{
				Operation: Get,
				Key:       key,
				Value:     value,
			}, nil) {
				return
			}

		}
	}

}
