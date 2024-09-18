package repl

import (
	"bytes"
	"cathon/token"
	"fmt"
	"strings"
	"testing"
)

// A mock lexer function to use for the test, if not already implemented
func TestREPL(t *testing.T) {
	// Simulate user input
	input := `let five = 5;
	`

	tests := []struct {
		expectedType    string
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
	}

	in := strings.NewReader(input)
	var out bytes.Buffer
	Start(in, &out)
	rawOutput := strings.TrimLeft(out.String(), ">> ")
	realOutputs := strings.Split(rawOutput, "\n")

	for i, tt := range tests {

		var expectedOut bytes.Buffer
		fmt.Fprintf(&expectedOut, "{Type:%s Literal:%s}", tt.expectedType, tt.expectedLiteral)

		if realOutputs[i] != expectedOut.String() {
			t.Fatalf("tests[%d] - Wrong output. expected=%q, got%q", i, expectedOut.String(), realOutputs[i])
		}
	}
}
