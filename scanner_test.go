package esprimago

import (
	"fmt"
	"testing"
)

func test_scanning_multi_line_comment(t *testing.T) {
	scanner := NewScanner("var foo=1; /* \"330413500\" */", &ParserOptions{Comment: true})

	results := make([]string, 0)

	var token *Token

	for token == nil || token.Type != EOF {
		for _, comment := range scanner.ScanComments() {
			results = append(results, fmt.Sprintf("{%v}-{%v}", comment.Start, comment.End))
		}

		token = scanner.Lex()
	}

	if results[0] != "11-28" {
		t.Errorf("failed scanning multi line coment expected start-end: %v", results[0])
	}
}
