package esprimago

import "./char"

type Scanner struct {
	Source     []rune
	Index      int
	LineNumber int
	LineStart  int

	length  int
	options *ParserOptions
}

type ParserOptions struct {
	Comment      bool
	ErrorHandler *ErrorHandler
}

func NewScanner(code string, options *ParserOptions) *Scanner {
	return &Scanner{
		Source: []rune(code),
		Index:  0,

		length:  len(code),
		options: options,
	}
}

func (s *Scanner) eof() bool {
	return s.Index >= s.length
}

func (s *Scanner) nextChar() rune {
	return s.nextCharOffset(0)
}

func (s *Scanner) nextCharOffset(offset int) rune {
	return s.Source[s.Index+offset]
}

func (s *Scanner) skipSingleLineComment(offset int) []*Comment {
	comments := make([]*Comment, 0)
	start := 0
	loc := &SourceLocation{}

	if s.options.Comment {
		start = s.Index - offset
		loc.Start = &Position{Line: s.LineNumber, Column: s.Index - s.LineStart - offset}
		loc.End = &Position{}
	}

	for !s.eof() {
		ch := s.nextChar()
		s.Index += 1

		if char.IsLineTerminator(ch) {
			if s.options.Comment {
				loc.End.Line = s.LineNumber
				loc.End.Column = s.Index - s.LineStart - 1

				comment := &Comment{
					MultiLine: false,
					Slice:     []int{start + offset, s.Index - 1},
					Start:     start,
					End:       s.Index - 1,
					Loc:       loc,
				}

				comments = append(comments, comment)
			}

			if ch == char.CARIAGE_RETURN && s.nextChar() == char.NEWLINE_FEED {
				s.Index += 1
			}

			s.LineNumber += 1
			s.LineStart = s.Index

			return comments
		}
	}

	if s.options.Comment {
		loc.End.Line = s.LineNumber
		loc.End.Column = s.Index - s.LineStart

		comment := &Comment{
			MultiLine: false,
			Slice:     []int{start + offset, s.Index},
			Start:     start,
			End:       s.Index,
			Loc:       loc,
		}

		comments = append(comments, comment)
	}

	return comments
}

func (s *Scanner) skipMultiLineComment() []*Comment {
	comments := make([]*Comment, 0)
	start := 0
	loc := &SourceLocation{}

	if s.options.Comment {
		start = s.Index - 2
		loc.Start = &Position{Line: s.LineNumber, Column: s.Index - s.LineStart - 2}
	}

	for !s.eof() {

		ch := s.nextChar()

		if char.IsLineTerminator(ch) {

			if ch == char.CARIAGE_RETURN && s.nextCharOffset(1) == char.NEWLINE_FEED {
				s.Index += 1
			}
			s.LineNumber += 1
			s.Index += 1
			s.LineStart = s.Index

		} else if ch == char.ASTERICK {

			// Block comment ends with '*/'

			if s.nextCharOffset(1) == char.FORWARD_SLASH {

				s.Index += 2

				if s.options.Comment {
					loc.End = &Position{Line: s.LineNumber, Column: s.Index - s.LineStart}
					comment := &Comment{
						MultiLine: true,
						Slice:     []int{start + 2, s.Index - 2},
						Start:     start,
						End:       s.Index,
						Loc:       loc,
					}
					comments = append(comments, comment)
				}

				return comments
			}

			s.Index += 1

		} else {

			s.Index += 1

		}
	}

	if s.options.Comment {
		loc.End = &Position{s.LineNumber, s.Index - s.LineStart}
		comment := &Comment{
			MultiLine: true,
			Slice:     []int{start + 2, s.Index},
			Start:     start,
			End:       s.Index,
			Loc:       loc,
		}

		comments = append(comments, comment)
	}

	s.tolerateUnexpectedToken(nil)

	return comments
}

func (s *Scanner) tolerateUnexpectedToken(message *string) {
	if message == nil {
		*message = Messages_UnexpectedTokenIllegal
	}

	(*s.options.ErrorHandler).TolerateError(s.Index, s.LineNumber, s.Index-s.LineStart+1, *message)
}

func (s *Scanner) ScanComments() []*Comment {
	comments := make([]*Comment, 0)

	start := s.Index == 0
	for !s.eof() {
		ch := s.nextChar()

		if char.IsWhiteSpace(ch) {

			s.Index += 1

		} else if char.IsLineTerminator(ch) {
			// "\n" or "\r\n"

			s.Index += 1

			if ch == char.CARIAGE_RETURN && s.nextChar() == char.NEWLINE_FEED {
				s.Index += 1
			}

			s.LineNumber += 1
			s.LineStart = s.Index
			start = true

		} else if ch == char.FORWARD_SLASH {
			// start of "//" or "/*"

			ch = s.nextCharOffset(1)

			if ch == char.FORWARD_SLASH {

				s.Index += 2

				singleLineComments := s.skipSingleLineComment(2)

				if s.options.Comment {
					for _, comment := range singleLineComments {
						comments = append(comments, comment)
					}
				}

			} else if ch == char.ASTERICK {

				s.Index += 2

				multiLineComments := s.skipMultiLineComment()

				if s.options.Comment {
					for _, comment := range multiLineComments {
						comments = append(comments, comment)
					}
				}

			} else {
				break
			}
		} else if start && ch == char.MINUS_SIGN {
			// when start of '-' and '>'

			if s.nextCharOffset(1) == char.MINUS_SIGN && s.nextCharOffset(2) == char.GREATER_THAN_SIGN {
				// '-->' is a single-line comment
				s.Index += 3

				singleLineComments := s.skipSingleLineComment(3)

				if s.options.Comment {
					for _, comment := range singleLineComments {
						comments = append(comments, comment)
					}
				}

			} else {
				break
			}

		} else if ch == char.LESS_THAN_SIGN { /* && !IsModule NOTE: IsModule is being set in the C# port therefore skip adding here */

			// check '<' followed by '!--'
			if s.nextCharOffset(1) == char.EXCLAMATION &&
				s.nextCharOffset(2) == char.MINUS_SIGN &&
				s.nextCharOffset(3) == char.MINUS_SIGN {

				s.Index += 4 // '<!--'

				singleLineComments := s.skipSingleLineComment(4)

				if s.options.Comment {
					for _, comment := range singleLineComments {
						comments = append(comments, comment)
					}
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return comments
}

func (s *Scanner) Lex() *Token {
	if s.eof() {
		return &Token{
			Type:       EOF,
			LineNumber: s.LineNumber,
			LineStart:  s.LineStart,
			Start:      s.Index,
			End:        s.Index,
		}
	}

	return nil
}
