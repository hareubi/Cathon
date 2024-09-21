package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x =5;
	let y = 10;
	let foobar = 838383;
	`

	testLexer := lexer.New(input)
	testParser := New(testLexer)

	program := testParser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 Statements. got %d Statements", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, testStatement ast.Statement, name string) bool {
	if testStatement.TokenLiteral() != "let" {
		t.Errorf("testStatement.TokenLiteral is not 'let'. got %q", testStatement.TokenLiteral())
		return false
	}

	letStmt, ok := testStatement.(*ast.LetStatement)
	if !ok {
		t.Errorf("testStatement is not *ast.LetStatement. got %T", testStatement)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.TokenLiteral() != "let" {
		t.Errorf("letStmt.Name.TokenLiteral() is not '%s'. got %q", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
