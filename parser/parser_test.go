package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	testLexer := lexer.New(input)
	testParser := New(testLexer)

	program := testParser.ParseProgram()
	CheckParserErrors(t, &testParser)
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
		if !CheckLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}
func CheckLetStatement(t *testing.T, testStatement ast.Statement, name string) bool {
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

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 99993;
	`

	testLexer := lexer.New(input)
	testParser := New(testLexer)

	program := testParser.ParseProgram()
	CheckParserErrors(t, &testParser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 Statements. got %d Statements", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"5"},
		{"10"},
		{"99993"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !CheckReturnStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}
func CheckReturnStatement(t *testing.T, testStatement ast.Statement, name string) bool {
	if testStatement.TokenLiteral() != "return" {
		t.Errorf("testStatement.TokenLiteral is not 'return'. got %q", testStatement.TokenLiteral())
		return false
	}

	returnStmt, ok := testStatement.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("testStatement is not *ast.ReturnStatement. got %T", testStatement)
		return false
	}
	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.Name.TokenLiteral() is not 'return'. got %q", returnStmt.TokenLiteral())
		return false
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := `
	foobar
	`

	testLexer := lexer.New(input)
	testParser := New(testLexer)

	program := testParser.ParseProgram()
	CheckParserErrors(t, &testParser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 Statements. got %d Statements", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !CheckIdentifierExpression(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}
func CheckIdentifierExpression(t *testing.T, testStatement ast.Statement, name string) bool {

	ExpressionStmt, ok := testStatement.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("ExpressionStatement is not *ast.ExpressionStatement. got %T", testStatement)
		return false
	}

	ident, ok := ExpressionStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got=%T", ExpressionStmt.Expression)
		return false
	}
	if ident.TokenLiteral() != name {
		t.Errorf("ident.TokenLiteral() is not %q. got %q", name, ident.TokenLiteral())
		return false
	}
	if ident.Value != name {
		t.Errorf("ident.Value is not %q. got %q", name, ident.Value)
		return false
	}

	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	testLexer := lexer.New(input)
	testParser := New(testLexer)

	program := testParser.ParseProgram()
	CheckParserErrors(t, &testParser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 Statements. got %d Statements", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"5"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !CheckIntegerLiteralExpression(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}
func CheckIntegerLiteralExpression(t *testing.T, testStatement ast.Statement, name string) bool {

	ExpressionStmt, ok := testStatement.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("ExpressionStatement is not *ast.ExpressionStatement. got %T", testStatement)
		return false
	}

	literal, ok := ExpressionStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral. got=%T", ExpressionStmt.Expression)
		return false
	}
	if literal.TokenLiteral() != name {
		t.Errorf("literal.TokenLiteral() is not %q. got %q", name, literal.TokenLiteral())
		return false
	}

	return true
}

func CheckParserErrors(t *testing.T, ParserP *Parser) {
	errors := ParserP.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error:%q", msg)
	}
	t.FailNow()
}
