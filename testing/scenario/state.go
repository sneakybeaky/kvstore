package scenario

import (
	"strings"
)

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*Lexer) stateFn

func startState(l *Lexer) stateFn {
	l.SkipWhitespace()

	if l.isEOF() {
		l.emit(ItemEOF)
		return startState
	}

	if strings.HasPrefix(l.inputToEnd(), put) {
		return lexPut
	} else if strings.HasPrefix(l.inputToEnd(), get) {
		return lexGet
	} else {
		return l.errorf("Expected GET or PUT but got %q", l.inputToEnd())
	}
}

func lexGet(l *Lexer) stateFn {
	l.pos += len(get)
	l.emit(ItemGet)
	return lexKey
}

func lexPut(l *Lexer) stateFn {
	l.pos += len(put)
	l.emit(ItemPut)
	return lexKey
}

func lexKey(l *Lexer) stateFn {
	l.SkipWhitespace()

	for {

		if strings.HasPrefix(l.inputToEnd(), " ") {
			l.emit(ItemKey)
			return lexValue
		}

		l.inc()

		if l.isEOF() {
			return l.errorf("Unexpected end of input")
		}
	}
}

func lexValue(l *Lexer) stateFn {
	l.SkipWhitespace()

	if strings.HasPrefix(l.inputToEnd(), notFound) {
		l.pos += len(notFound)
		l.emit(ItemNotFound)
		return startState
	}

	for {

		if strings.HasPrefix(l.inputToEnd(), newline) {
			l.emit(ItemValue)
			return startState
		}

		l.inc()

		if l.isEOF() {
			return l.errorf("Unexpected end of input")
		}
	}
}
