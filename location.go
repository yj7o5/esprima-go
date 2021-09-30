package esprimago

import "fmt"

type Location struct {
	Start  *Position
	End    *Position
	Source *string
}

func (l *Location) String() string {
	suffix := ""
	if l.Source != nil {
		suffix = fmt.Sprintf(": %v", l.Source)
	}

	return fmt.Sprintf("{%v}...{%v}{%v}", l.Start, l.End, suffix)
}

// TODO: Implement hash equality functions

// TODO: Implement Equals, Deconstruct functions

type SourceLocation struct {
	Start *Position
	End   *Position
}
