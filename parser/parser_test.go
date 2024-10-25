package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"fmt"
	"strconv"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	expectedIdentifiers := []string{"x", "y", "foobar"}
	CheckParseStatements(t, input, 3, expectedIdentifiers, CheckLetStatement)
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 99993;
	`
	expectedIdentifiers := []string{"5", "10", "99993"}
	CheckParseStatements(t, input, 3, expectedIdentifiers, CheckReturnStatement)
}

func TestIdentifierExpression(t *testing.T) {
	input := `
	foobar;
	mewo;
	`
	expectedIdentifiers := []string{"foobar", "mewo"}
	CheckParseStatements(t, input, 2, expectedIdentifiers, CheckIdentifierExpression)
}

func TestBoolExpression(t *testing.T) {
	input := `
 	false;
	true;
	false;
	`
	expectedIdentifiers := []string{"false", "true", "false"}

	CheckParseStatements(t, input, 3, expectedIdentifiers, CheckBoolExpression)
}

func CheckParseStatements(t *testing.T, input string, expectedStmtCount int, expectedIdentifiers []string, checkFunc func(*testing.T, ast.Statement, string)) {
	testLexer := lexer.New(input)
	testParser := New(testLexer)

	program := testParser.ParseProgram()
	CheckParserErrors(t, testParser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != expectedStmtCount {
		t.Fatalf("program.Statements does not contain %d Statements. got %d Statements", expectedStmtCount, len(program.Statements))
	}

	for i, expectedIdentifier := range expectedIdentifiers {
		stmt := program.Statements[i]
		checkFunc(t, stmt, expectedIdentifier)
	}
}

func CheckLetStatement(t *testing.T, testStatement ast.Statement, name string) {
	if testStatement.TokenLiteral() != "let" {
		t.Errorf("testStatement.TokenLiteral is not 'let'. got %q", testStatement.TokenLiteral())
	}
	letStmt, ok := testStatement.(*ast.LetStatement)
	if !ok {
		t.Fatalf("testStatement is not *ast.LetStatement. got %T", testStatement)
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s'. got=%s", name, letStmt.Name.Value)
	}
	if letStmt.TokenLiteral() != "let" {
		t.Errorf("letStmt.Name.TokenLiteral() is not '%s'. got %q", name, letStmt.Name.TokenLiteral())
	}
}

func CheckReturnStatement(t *testing.T, testStatement ast.Statement, name string) {
	if testStatement.TokenLiteral() != "return" {
		t.Errorf("testStatement.TokenLiteral is not 'return'. got %q", testStatement.TokenLiteral())
	}
	returnStmt, ok := testStatement.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("testStatement is not *ast.ReturnStatement. got %T", testStatement)
	}
	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.Name.TokenLiteral() is not 'return'. got %q", returnStmt.TokenLiteral())
	}
}

func CheckIdentifierExpression(t *testing.T, testStatement ast.Statement, name string) {
	ExpressionStmt, ok := testStatement.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("ExpressionStatement is not *ast.ExpressionStatement. got %T", testStatement)
	}
	ident, ok := ExpressionStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got=%T", ExpressionStmt.Expression)
	}
	if ident.TokenLiteral() != name {
		t.Errorf("ident.TokenLiteral() is not %q. got %q", name, ident.TokenLiteral())
	}
	if ident.Value != name {
		t.Errorf("ident.Value is not %q. got %q", name, ident.Value)
	}
}

func CheckBoolExpression(t *testing.T, testStatement ast.Statement, isTrue string) {
	ExpressionStmt, ok := testStatement.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("ExpressionStatement is not *ast.ExpressionStatement. got %T", testStatement)
	}
	ident, ok := ExpressionStmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got=%T", ExpressionStmt.Expression)
	}
	if ident.TokenLiteral() != isTrue {
		t.Errorf("ident.TokenLiteral() is not %q. got %q", isTrue, ident.TokenLiteral())
	}
	if ident.Value != (isTrue == "true") {
		t.Errorf("ident.Value is not %q. got %q", isTrue, strconv.FormatBool(ident.Value))
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `
	5;
	10
	`

	testLexer := lexer.New(input)
	testParser := New(testLexer)

	program := testParser.ParseProgram()
	CheckParserErrors(t, testParser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 Statements. got %d Statements", len(program.Statements))
	}

	test := "5"

	testStatement := program.Statements[0]
	ExpressionStmt, ok := testStatement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got %T", program.Statements[0])
	}

	literal, ok := ExpressionStmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral. got=%T", ExpressionStmt.Expression)
	}
	if literal.TokenLiteral() != test {
		t.Fatalf("literal.TokenLiteral() is not %q. got %q", test, literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}
	for _, tt := range prefixTests {
		testLexer := lexer.New(tt.input)
		testParser := New(testLexer)

		program := testParser.ParseProgram()
		CheckParserErrors(t, testParser)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 Statements. got %d Statements", len(program.Statements))
		}

		testStatement := program.Statements[0]
		ExpressionStmt, ok := testStatement.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("ExpressionStatement[0] is not *ast.ExpressionStatement. got %T", program.Statements[0])
		}

		exp, ok := ExpressionStmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expected *ast.PrefixExpression. got=%T", ExpressionStmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("literal.TokenLiteral() is not %q. got %q", tt.operator, exp.Operator)
		}
		CheckIntegerLiteral(t, exp.Right, tt.integerValue)
	}
}

func CheckIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return
	}
	if integ.Value != value {
		t.Errorf("integ.Value is not %d. got %d", value, integ.Value)
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.Value is not %d. got %s", value, integ.TokenLiteral())
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"100 + 1;", 100, "+", 1},
		{"101-2", 101, "-", 2},
		{"102 /3;", 102, "/", 3},
		{"103* 4", 103, "*", 4},
		{"104 < 5", 104, "<", 5},
		{"105 > 6", 105, ">", 6},
		{"106 == 7", 106, "==", 7},
		{"107 != 8", 107, "!=", 8},
	}
	for i, tt := range infixTests {
		testLexer := lexer.New(tt.input)
		testParser := New(testLexer)

		program := testParser.ParseProgram()
		CheckParserErrors(t, testParser)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Errorf("program%d.Statements does not contain 1 Statements. got %d Statements", i, len(program.Statements))
		}
		testStatement := program.Statements[0]
		ExpressionStmt, ok := testStatement.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("ExpressionStatement[0] is not *ast.ExpressionStatement. got %T", program.Statements[0])
		}

		exp, ok := ExpressionStmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expected *ast.InfixExpression. got=%T", ExpressionStmt.Expression)
		}

		CheckIntegerLiteral(t, exp.Left, tt.leftValue)
		if exp.Operator != tt.operator {
			t.Errorf("tt.operator is not %q. got %q", tt.operator, exp.Operator)
		}
		CheckIntegerLiteral(t, exp.Right, tt.rightValue)
	}
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
