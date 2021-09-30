package esprimago

import (
	"fmt"
	"reflect"
)

type Position struct {
	Line   int
	Column int
}

func NewPosition(line, col int) *Position {
	if line < 0 {
		panic(fmt.Sprintf("\"line\" %v is out of range", line))
	}

	if (line <= 0 || col < 0) && (line != 0 || col != 0) {
		panic(fmt.Sprintf("\"column\" %v is out of range", col))
	}

	return &Position{line, col}
}

func (p *Position) EqualInterface(o interface{}) bool {
	oType := reflect.TypeOf(o)
	rType := reflect.TypeOf(p)

	return oType.AssignableTo(rType)
}

func (p *Position) EqualPosition(o *Position) bool {
	return p.Line == o.Line && p.Column == o.Column
}

// TODO: Implement hash functions

// TODO: Implement equality functions
