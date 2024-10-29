package parser

import (
	"cathon/ast"
	"cathon/lexer"
	"fmt"
	"strconv"
	"strings"
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

func TestIntegerLiteralExpression(t *testing.T) {
	input := `
	5;
	10
	`
	expectedIdentifiers := []string{"5", "10"}
	CheckParseExpression(t, input, 2, expectedIdentifiers, CheckIntegerLiteralExpression)
}
func TestBoolExpression(t *testing.T) {
	input := `
 	false;
	true;
	false;
	`
	expectedIdentifiers := []string{"false", "true", "false"}

	CheckParseExpression(t, input, 3, expectedIdentifiers, CheckBoolExpression)
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
func TestInfixExpressions(t *testing.T) {
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
func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	CheckInfixExpression(t, exp.Condition, "x", "<", "y")

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	CheckIdentifierExpression(t, consequence.Expression, "x")

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	CheckInfixExpression(t, exp.Condition, "x", "<", "y")

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	CheckIdentifierExpression(t, consequence.Expression, "x")

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}
	CheckIdentifierExpression(t, alternative.Expression, "y")
}
func TestIdentifierExpression(t *testing.T) {
	input := `
	foobar;
	mewo;
	`
	expectedIdentifiers := []string{"foobar", "mewo"}
	CheckParseExpression(t, input, 2, expectedIdentifiers, CheckIdentifierExpression)
}
func TestOperatorPrecedenceExpression(t *testing.T) {
	input := `
	1 + (2 +3)+ 4
	1 + 2 +(3 + 4)
	`
	CheckParseExpression(t, input, 2, []string{"((1 + (2 + 3)) + 4)","((1 + 2) + (3 + 4))"}, CheckOperatorPrecedenceExpression)
}
func TestFunctionExpression(t *testing.T) {
	input := `
	fn(x, y) { x + y; }
	`
	CheckParseExpression(t, input, 1, []string{"x,y,x,+,y"}, CheckFunctionExpression)
}
func TestCallExpression(t *testing.T) {
	input := `
	add(1,2 * 3,4+ 5);
	`
	CheckParseExpression(t, input, 1, []string{"add,1,2,*,3,4,+,5"}, CheckCallExpression)
}

func CheckOperatorPrecedenceExpression(t *testing.T, exp ast.Expression, expected  string) {
			actual := exp.String()
		if actual != expected {
			t.Errorf("expected=%q, got=%q", expected , actual)
		}
}
func CheckParseExpression(t *testing.T, input string, expectedStmtCount int, expectedIdentifiers []string, checkFunc func(*testing.T, ast.Expression, string)) {
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
		exprStmt, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("ExpressionStatement is not *ast.ExpressionStatement. got %T", program.Statements[i])
		}
		expr := exprStmt.Expression
		checkFunc(t, expr, expectedIdentifier)
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
func CheckIntegerLiteralExpression(t *testing.T, testExpression ast.Expression, name string) {
	ident, ok := testExpression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got=%T", testExpression)
	}
	if ident.TokenLiteral() != name {
		t.Errorf("ident.TokenLiteral() is not %q. got %q", name, ident.TokenLiteral())
	}
	if strconv.FormatInt(ident.Value, 10) != name {
		t.Errorf("ident.Value is not %q. got %q", name, ident.Value)
	}
}
func CheckBoolExpression(t *testing.T, testExpression ast.Expression, isTrue string) {
	ident, ok := testExpression.(*ast.Boolean)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got=%T", testExpression)
	}
	if ident.TokenLiteral() != isTrue {
		t.Errorf("ident.TokenLiteral() is not %q. got %q", isTrue, ident.TokenLiteral())
	}
	if ident.Value != (isTrue == "true") {
		t.Errorf("ident.Value is not %q. got %q", isTrue, strconv.FormatBool(ident.Value))
	}
}
func CheckIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("il not *ast.IntegerLiteral. got=%T", il)
	}
	if integ.Value != value {
		t.Errorf("integ.Value is not %d. got %d", value, integ.Value)
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.Value is not %d. got %s", value, integ.TokenLiteral())
	}
}
func CheckIdentifierExpression(t *testing.T, testExpression ast.Expression, name string) {
	ident, ok := testExpression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier. got=%T", testExpression)
	}
	if ident.TokenLiteral() != name {
		t.Errorf("ident.TokenLiteral() is not %q. got %q", name, ident.TokenLiteral())
	}
	if ident.Value != name {
		t.Errorf("ident.Value is not %q. got %q", name, ident.Value)
	}
}
func CheckInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
	}

	CheckLiteralExpression(t, opExp.Left, left)

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
	}

	CheckLiteralExpression(t, opExp.Right, right)
}
func CheckLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		CheckIntegerLiteral(t, exp, int64(v))
		return
	case int64:
		CheckIntegerLiteral(t, exp, v)
		return
	case string:
		CheckIdentifierExpression(t, exp, v)
		return
	case bool:
		CheckBoolExpression(t, exp, strconv.FormatBool(v))
		return
	}
	t.Errorf("type of exp not handled. got=%T", exp)
}
func CheckFunctionExpression(t *testing.T, exp ast.Expression, expected string) {
	
	expectedArray := strings.Split(expected, ",")
	for i ,value := range(expectedArray) {
		expectedArray[i] =strings.NewReplacer(",", "").Replace(value)
	}
	function, ok := exp.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expected *ast.FunctionLiteral. got=%T", exp)
	}
	CheckLiteralExpression(t, function.Parameters[0], expectedArray[0])
	CheckLiteralExpression(t, function.Parameters[1], expectedArray[1])

	bodyStmt, ok:= function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	} 
	CheckInfixExpression(t, bodyStmt.Expression, expectedArray[2], expectedArray[3], expectedArray[4])
}
func CheckCallExpression(t *testing.T, exp ast.Expression, expected string) {
	
	expectedArray := strings.Split(expected, ",")
	for i ,value := range(expectedArray) {
		expectedArray[i] =strings.NewReplacer(",", "").Replace(value)
	}
	call, ok := exp.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected *ast.CallExpression. got=%T", exp)
	}
	CheckIdentifierExpression(t, call.Function, expectedArray[0])
	if len(call.Arguments) != 3{
		t.Fatalf("expected 3 arguments, got %d",len(call.Arguments))
	}
	data, err:= strconv.ParseInt(expectedArray[1],10,16)
	if err!= nil {
		t.Fatalf("error parsing testdata. got=%q",err.Error())
	}
	CheckLiteralExpression(t, call.Arguments[0], data)
		data1, err:= strconv.ParseInt(expectedArray[2],10,16)
	if err!= nil {
		t.Fatalf("error parsing testdata. got=%q",err.Error())
	}
		data2, err:= strconv.ParseInt(expectedArray[4],10,16)
	if err!= nil {
		t.Fatalf("error parsing testdata. got=%q",err.Error())
	}
		data3, err:= strconv.ParseInt(expectedArray[5],10,16)
	if err!= nil {
		t.Fatalf("error parsing testdata. got=%q",err.Error())
	}
		data4, err:= strconv.ParseInt(expectedArray[7],10,16)
	if err!= nil {
		t.Fatalf("error parsing testdata. got=%q",err.Error())
	}
	CheckInfixExpression(t, call.Arguments[1], data1 ,expectedArray[3], data2)
	CheckInfixExpression(t, call.Arguments[2], data3, expectedArray[6], data4)
}