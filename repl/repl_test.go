package repl

import (
	"bytes"
	"strings"
	"testing"
)

// A mock lexer function to use for the test, if not already implemented
func TestREPL(t *testing.T) {
	// Simulate user input
	input := `2 +3;
	`

	tests := []string{"5"}

	in := strings.NewReader(input)
	var out bytes.Buffer
	Start(in, &out)
	rawOutput := strings.TrimLeft(out.String(), "> ")
	realOutputs := strings.Split(rawOutput, "\n")

	for i, tt := range tests {

		if realOutputs[i] != tt {
			t.Fatalf("tests[%d] - Wrong output. expected=%q, got%q", i, tt, realOutputs[i])
		}
	}
}
