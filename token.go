package esprimago

import "regexp"

type TokenType int

const (
	BooleanLiteral TokenType = iota
	EOF
	Identifier
	Keyword
	NullLiteral
	NumericLiteral
	Punctuator
	StringLiteral
	RegularExpression
	Template
)

type Token struct {
	Type       TokenType
	Literal    *string
	Start      int
	End        int
	LineNumber int
	LineStart  int

	Location *Location

	// Numeric Literals
	Octal                 bool
	NotEscapeSequenceHead *rune

	// Templates
	Head         bool
	Tail         bool
	RawTemplate  *string
	BooleanValue bool
	NumericValue float64
	Value        *interface{}
	RegexValue   *regexp.Regexp
}

func (t *Token) Clear() {
	t.Type = BooleanLiteral
	t.Literal = nil
	t.Start = 0
	t.End = 0
	t.LineNumber = 0
	t.LineStart = 0
	t.Location = nil
	t.Octal = false
	t.Head = false
	t.Tail = false
	t.RawTemplate = nil
	t.BooleanValue = false
	t.NumericValue = 0
	t.Value = nil
	t.RegexValue = nil
}
