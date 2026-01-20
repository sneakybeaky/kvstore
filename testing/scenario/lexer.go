package scenario

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const EOF rune = 0

// Lexer holds the state of the scanner.
type Lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this Item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan Item // channel of scanned items.
	state stateFn
}

func NewLexer(name, input string) *Lexer {
	return &Lexer{state: startState, name: name, input: input, items: make(chan Item, 2)}
}

func (l *Lexer) NextItem() Item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}

}

func (l *Lexer) skipWhitespace() {
	for {
		ch := l.next()

		if ch == EOF {
			l.emit(ItemEOF)
			break
		}

		if !unicode.IsSpace(ch) {
			l.dec()
			break
		}

	}
}

func (l *Lexer) emit(it ItemType) {
	l.items <- Item{Typ: it, Val: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) next() rune {
	if l.pos >= utf8.RuneCountInString(l.input) {
		l.width = 0
		return EOF
	}

	result, width := utf8.DecodeRuneInString(l.input[l.pos:])

	l.width = width
	l.pos += l.width
	return result
}

func (l *Lexer) inputToEnd() string {
	return l.input[l.pos:]
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{
		Typ: ItemError,
		Val: fmt.Sprintf(format, args...),
	}

	return nil
}

func (l *Lexer) inc() {
	l.pos++
	if l.pos >= utf8.RuneCountInString(l.input) {
		l.emit(ItemEOF)
	}
}

func (l *Lexer) dec() {
	l.pos--
}

func (l *Lexer) isEOF() bool {
	return l.pos >= len(l.input)
}

func (l *Lexer) Shutdown() {
	close(l.items)
}
