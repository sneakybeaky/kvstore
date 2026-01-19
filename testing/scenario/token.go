package scenario

import "fmt"

type ItemType int

// ItemType identifies the type of lex items.
const (
	ItemError ItemType = iota // error occurred;
	ItemEOF
	ItemPut
	ItemGet
	ItemKey
	ItemValue
)

const put string = "PUT"
const get string = "GET"
const newline string = "\n"

// Item represents a token returned from the scanner.
type Item struct {
	Typ ItemType // Type, such as itemNumber.
	Val string   // Value, such as "23.2".
}

func (i Item) String() string {
	switch i.Typ {
	case ItemEOF:
		return "EOF"
	case ItemError:
		return i.Val
	}
	if len(i.Val) > 10 {
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}
