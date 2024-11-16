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
	let y = true;
	let foobar = x;
	`
	expectedIdentifiers := []string{"x,5", "y,true", "foobar,x"}
	CheckParseStatements(t, input, 3, expectedIdentifiers, CheckLetStatement)
}
func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return true;
	return x;
	`
	expectedIdentifiers := []string{"5", "true", "x"}
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
func CheckLetStatement(t *testing.T, testStatement ast.Statement, expected string) {
		expectedArray := strings.Split(expected, ",")
	for i ,value := range(expectedArray) {
		expectedArray[i] =strings.NewReplacer(",", "").Replace(value)
	}
	if testStatement.TokenLiteral() != "let" {
		t.Errorf("testStatement.TokenLiteral is not 'let'. got %q", testStatement.TokenLiteral())
	}
	letStmt, ok := testStatement.(*ast.LetStatement)
	if !ok {
		t.Fatalf("testStatement is not *ast.LetStatement. got %T", testStatement)
	}
	if letStmt.Name.Value != expectedArray[0] {
		t.Errorf("letStmt.Name.Value is not '%s'. got=%s", expectedArray[0], letStmt.Name.Value)
	}
		if letStmt.Value.String() != expectedArray[1] {
		t.Errorf("letStmt.Value.String() is not '%s'. got=%s", expectedArray[1], letStmt.Value.String())
	}
	if letStmt.TokenLiteral() != "let" {
		t.Errorf("letStmt.Name.TokenLiteral() is not '%s'. got %q", expected, letStmt.Name.TokenLiteral())
	}
}
func CheckReturnStatement(t *testing.T, testStatement ast.Statement, expected string) {
	if testStatement.TokenLiteral() != "return" {
		t.Errorf("testStatement.TokenLiteral is not 'return'. got %q", testStatement.TokenLiteral())
	}
	returnStmt, ok := testStatement.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("testStatement is not *ast.ReturnStatement. got %T", testStatement)
	}
	if returnStmt.ReturnValue.String() != expected {
		t.Errorf("returnStmt.ReturnValue.String() is not 'return'. got %q", expected)
	}
	if returnStmt.TokenLiteral() != "return" {
		t.Errorf("returnStmt.TokenLiteral() is not 'return'. got %q", returnStmt.TokenLiteral())
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

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingEmptyArrayLiterals(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 0 {
		t.Errorf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	CheckIntegerLiteral(t, array.Elements[0], 1)
	CheckInfixExpression(t, array.Elements[1], 2, "*", 2)
	CheckInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	CheckIdentifierExpression(t, indexExp.Left, "myArray")
	CheckInfixExpression(t, indexExp.Index, 1, "+", 1) 
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[literal.String()]
		CheckIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		boolean, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.BooleanLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[boolean.String()]
		CheckIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		integer, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[integer.String()]

		CheckIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			CheckInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			CheckInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			CheckInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}
