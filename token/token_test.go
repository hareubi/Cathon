package token

import "testing"

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"fn", FUNCTION},
		{"let", LET},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"foobar", IDENT}, // Non-keyword, should return IDENT
		{"x", IDENT},      // Single character identifier
	}

	for i, tt := range tests {
		tokType := LookupIdent(tt.input)
		if tokType != tt.expected {
			t.Errorf("expected %q for input %q, got %q", tt.expected, tt.input, tokType)
		}
	}
}
